package apptianlang

import (
	"TianlangCapturer/src/api"
	"TianlangCapturer/src/model"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	mDownloadTask   chan *model.VideoInfo
	mHistoryRecord  chan *model.VideoInfo
	mKnownVideoInfo sync.Map
	mCountLimit     int
	mCountMutex     sync.Mutex
)

// InitDownloader InitDownloader
func InitDownloader() {
	api.Info("Create download worker")
	mDownloadTask = make(chan *model.VideoInfo, 4096)
	for workerNum := runtime.NumCPU() * 2; workerNum > 0; workerNum-- {
		go workerGoroutine()
	}

	api.Info("Load history record..")
	if recordFile, err := os.Open("record.dat"); err == nil {
		defer recordFile.Close()
		reader := bufio.NewReader(recordFile)
		for {
			lineStr, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			if data, err := base64.StdEncoding.DecodeString(lineStr); err == nil {
				videoInfo := &model.VideoInfo{}
				if model.Unmarshal(data, videoInfo) == nil && videoInfo.ID > 0 {
					mKnownVideoInfo.Store(videoInfo.GetMD5(), videoInfo)
				}
			}
		}
	}

	mHistoryRecord = make(chan *model.VideoInfo, 128)
	go func() {
		var (
			err  error
			file *os.File
		)
		for videoInfo := range mHistoryRecord {
			api.Info("[Doen]Task:%s", videoInfo.GetMD5())
			data, _ := model.Marshal(videoInfo)
			base64Data := base64.StdEncoding.EncodeToString(data)
			if file, err = os.OpenFile("record.dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err == nil {
				file.WriteString(base64Data)
				file.WriteString("\r\n")
				file.Close()
			} else {
				api.Error("[Record]Open 'record.dat' error:%s recordData:\n%s", err.Error(), base64Data)
			}
		}
	}()
}

// XORAllFileData XORAllFileData
func XORAllFileData(rootDir string, key []byte) {
	api.Debug("[XORAllFileData]%s", rootDir)
	fileInfos, _ := ioutil.ReadDir(rootDir)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			XORAllFileData(filepath.Join(rootDir, fileInfo.Name()), key)
		} else {
			filePath := filepath.Join(rootDir, fileInfo.Name())
			if file, err := os.OpenFile(filePath, os.O_RDWR, os.ModePerm); err == nil {
				keyLen := len(key)
				buff := make([]byte, 1024)
				keyIndex := 0
				for {
					n, err := file.Read(buff)
					if err != nil {
						break
					}
					data := buff[:n]
					for i, d := range data {
						data[i] = d ^ key[keyIndex%keyLen]
						keyIndex++
					}
					file.Seek(int64(-n), os.SEEK_CUR)
					file.Write(data)
				}
				file.Close()
			}
		}
	}
}

func isCanAddDownloadTask() bool {
	if mConfigInfo.MaxLimit > 0 {
		mCountMutex.Lock()
		if mCountLimit >= mConfigInfo.MaxLimit {
			mCountMutex.Unlock()
			return false
		}
		mCountLimit++
		mCountMutex.Unlock()
	}
	return true
}

// AddDownloadTask AddDownloadTask
func AddDownloadTask(videoInfo *model.VideoInfo) {
	if isCanAddDownloadTask() {
		videoMD5 := videoInfo.GetMD5()
		if _, loaded := mKnownVideoInfo.LoadOrStore(videoMD5, videoInfo); loaded {
			return
		}
		mKnownVideoInfo.Store(videoMD5, videoInfo)
		api.Info("[AddTask]ID:%d MD5:%s", videoInfo.ID, videoInfo.GetMD5())
		mDownloadTask <- videoInfo
	}
}

func workerGoroutine() {
	for videoInfo := range mDownloadTask {
		if workerExecTask(videoInfo) == false {
			if videoInfo.RetryCount < 10 {
				videoInfo.RetryCount++
				mDownloadTask <- videoInfo
			}
		} else {
			//添加已完成记录
			mHistoryRecord <- videoInfo
		}
	}
}

func workerExecTask(videoInfo *model.VideoInfo) bool {
	defer api.ProtectError()

	//HTTP查询资源是否存在(去重)
	var err error
	var resp *http.Response
	videoMD5 := videoInfo.GetMD5()
	if resp, err = http.Post(mConfigInfo.WebAPI+"video/query/exist", "text/plain", bytes.NewReader([]byte(videoInfo.Title))); err == nil {
		data, _ := ioutil.ReadAll(resp.Body)
		api.Info("[WebAPI]%s StateCode:%d Details:\n%s", videoMD5, resp.StatusCode, string(data))
		if string(data) == "exists" {
			return true
		}
	} else {
		api.Error("[WebAPI]%s Error:%s", videoMD5, err.Error())
		return false
	}

	//创建主目录
	mRootDir := fmt.Sprintf("Data/%s", videoMD5)
	os.MkdirAll(mRootDir, os.ModePerm)

	//下载视频文件
	downloader := NewDownloaderM3U8(videoInfo.M3U8Path, mRootDir)
	if err = downloader.WaitDownload(); err != nil {
		api.Error("[Download]%s Error:%s", videoMD5, err.Error())
		return false
	}

	// //加密目录下所有文件
	// api.XORAllFileData(mRootDir, api.DefaultXORKey)

	//SCP上传到资源服务器
	var msg string
	if msg, err = api.SCP(mConfigInfo.SCPAddr, mConfigInfo.SCPPort, mConfigInfo.SCPUser, mConfigInfo.SCPPass, mRootDir, mConfigInfo.SavePath); err == nil {
		api.Info("[SCP]%s Done", videoMD5)
	} else {
		api.Error("[SCP]%s Error:%s\ndetails:\n%s", videoMD5, err.Error(), msg)
		return false
	}

	//HTTP通知资源接口(入库)
	type UploadPostJSON struct {
		VideoName      string `json:"videoName"`
		MD5Name        string `json:"md5Name"`
		VideoType      string `json:"videoType"`
		VideoChildType string `json:"videoChildType"`
		Tags           string `json:"tags"`
	}
	var jsonRawData UploadPostJSON
	jsonRawData.VideoName = videoInfo.Title
	jsonRawData.MD5Name = videoMD5
	jsonRawData.VideoType = videoInfo.MajorCategory
	jsonRawData.VideoChildType = videoInfo.MinorCategory
	jsonRawData.Tags = videoInfo.Tags
	jsonByteData, _ := json.Marshal(&jsonRawData)
	if resp, err = http.Post(mConfigInfo.WebAPI+"video/upload", "application/json", bytes.NewReader(jsonByteData)); err == nil {
		data, _ := ioutil.ReadAll(resp.Body)
		api.Info("[WebAPI]%s StateCode:%d Details:\n%s", videoMD5, resp.StatusCode, string(data))
	} else {
		api.Error("[WebAPI]%s Error:%s", videoMD5, err.Error())
		return false
	}

	//HTTP通知切片接口(处理缩略图以及封面)
	if resp, err = http.Get(fmt.Sprintf(mConfigInfo.WebAPI+"/video/videoHandler?videoPath=%s&md5Name=%s", url.QueryEscape(mConfigInfo.SavePath), videoMD5)); err == nil {
		data, _ := ioutil.ReadAll(resp.Body)
		result := string(data)
		api.Info("[WebAPI]%s videoHandler StateCode:%d Details:\n%s", videoMD5, resp.StatusCode, result)
		if result != "success" {
			api.Error("[WebAPI]%s videoHandler result:%s", videoMD5, result)
			return false
		}
	} else {
		api.Error("[WebAPI]%s videoHandler Error:%s", videoMD5, err.Error())
		return false
	}

	//清空目录
	api.Info("[Clear]%s", mRootDir)
	os.RemoveAll(mRootDir)
	return true
}
