package appqianjiao

import (
	"TianlangCapturer/src/api"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// ConfigInfo ConfigInfo
type ConfigInfo struct {
	WebAPI      string `json:"webAPI"`
	SCPAddr     string `json:"scpAddr"`
	SCPPort     string `json:"scpPort"`
	SCPUser     string `json:"scpUser"`
	SCPPass     string `json:"scpPass"`
	SavePath    string `json:"savePath"`
	MaxLimit    int    `json:"maxLimit"`
	MechineCode string `json:"mechineCode"`
}

var (
	mConfigInfo ConfigInfo
)

// QianjiaoVideoModule QianjiaoVideoModule
type QianjiaoVideoModule struct {
	mRecorder     *api.Recorder
	mDownloadTask chan *VideoInfo
}

// Name Name
func (p *QianjiaoVideoModule) Name() string { return "Qianjiao.Video" }

// Run Run
func (p *QianjiaoVideoModule) Run() (err error) {

	// 加载配置
	mConfigInfo.WebAPI = "http://192.168.100.49:8805/"
	mConfigInfo.SCPAddr = "192.168.100.51"
	mConfigInfo.SCPPort = "22"
	mConfigInfo.SCPUser = "root"
	mConfigInfo.SCPPass = "Dd778899"
	mConfigInfo.SavePath = "~/scp-test"
	mConfigInfo.MaxLimit = 0
	mConfigInfo.MechineCode = "5c3b9e8cf47a413a"

	// 初始化记录器
	p.mRecorder = api.GetHistory(p, p.onHistoryHandleFunc)

	// 开始捕获
	for p.initDownloader(); p.loadConfig(); {
		disguiser := NewDisguiser(mConfigInfo.MechineCode)
		if err = disguiser.Login(); err != nil {
			api.Error("LoginError:%s", err.Error())
			return
		}

		// 获取所有专题
		var tagInfoList []TagInfo
		if tagInfoList, err = disguiser.GetTagByArea("专题"); err != nil {
			api.Error(err.Error())
			return
		}

		// 根据专题获取所有影片
		for _, tagInfo := range tagInfoList {
			api.Info("ID:%s Name:%s CoverURL:%s", tagInfo.ID, tagInfo.VideoTag, tagInfo.TagCover)

			for page := 0; true; page++ {

				// 获取该专题下某页影片
				var videoInfoList []VideoInfo
				if videoInfoList, err = disguiser.GetVideoByTag(tagInfo.ID, page); err != nil {
					return
				}

				// 填充分类标题
				for i := range videoInfoList {
					info := &videoInfoList[i]
					info.MajorCategory = "专题"
					info.MinorCategory = tagInfo.VideoTag
				}

				// 下载该页每部影片
				for _, videoInfo := range videoInfoList {
					// 该影片是否已下载
					if _, ok := p.mRecorder.FindRecord(videoInfo.ID); ok {
						api.Info("[Skip]ID:%s Name:%d URl:%s Headers:%s", videoInfo.ID, len(videoInfo.VideoTitle), videoInfo.VideoURL, videoInfo.Headers)
						continue
					}

					// 新建下载任务
					api.Info("[Task]ID:%s Name:%d URl:%s Headers:%s", videoInfo.ID, len(videoInfo.VideoTitle), videoInfo.VideoURL, videoInfo.Headers)
					pVideoInfo := &VideoInfo{
						ID:            videoInfo.ID,
						Tags:          videoInfo.Tags,
						VideoTitle:    videoInfo.VideoTitle,
						VideoCover:    videoInfo.VideoCover,
						VideoURL:      videoInfo.VideoURL,
						Headers:       videoInfo.Headers,
						MajorCategory: videoInfo.MajorCategory,
						MinorCategory: videoInfo.MinorCategory,
						RetryCount:    videoInfo.RetryCount,
						HeadersMap:    videoInfo.HeadersMap,
					}
					for isRetry := true; isRetry; {
						select {
						case p.mDownloadTask <- pVideoInfo:
							isRetry = false
						default:
							api.Warn("[Disguiser]The task queue is full!")
							time.Sleep(5 * time.Minute)
						}
					}
					time.Sleep(1 * time.Second)
				}

				if len(videoInfoList) < 24 {
					break
				}
			}
		}

		api.Warn("[Disguiser]Sleep(1h30m)")
		time.Sleep(1*time.Hour + 30*time.Minute)
	}
	return
}

func (p *QianjiaoVideoModule) loadConfig() bool {
	api.Info("[Config]ReloadUpdate")
	data, err := ioutil.ReadFile("QianjiaoCapturer.cfg")
	if err == nil {
		json.Unmarshal(data, &mConfigInfo)
	} else {
		data, _ = json.MarshalIndent(&mConfigInfo, "", "\t")
		ioutil.WriteFile("QianjiaoCapturer.cfg", data, os.ModePerm)
	}
	return true
}

func (p *QianjiaoVideoModule) onHistoryHandleFunc(data []byte) (key, value interface{}, err error) {
	result := new(VideoInfo)
	if err = json.Unmarshal(data, result); err != nil {
		return
	}
	key = result.ID
	value = result
	return
}

func (p *QianjiaoVideoModule) initDownloader() {
	api.Info("Create download worker")
	p.mDownloadTask = make(chan *VideoInfo, 4096)
	for workerNum := runtime.NumCPU() * 2; workerNum > 0; workerNum-- {
		go p.workerGoroutine()
	}
}

func (p *QianjiaoVideoModule) workerGoroutine() {
	defer api.ProtectError()
	defer api.Warn("[DownloadWorker]Exit")
	for videoInfo := range p.mDownloadTask {
		if err := p.workerExecTask(videoInfo); err != nil {
			api.Error(err.Error())
			if videoInfo.RetryCount < 10 {
				videoInfo.RetryCount++
				select {
				case p.mDownloadTask <- videoInfo:
				default:
					api.Warn("[DownloadWorker]The add retry task failed, task queue is full! task:%s", videoInfo.GetMD5())
				}
			} else {
				// 下载失败(尝试太多，不可能完成的任务)
				videoMD5 := videoInfo.GetMD5()
				mRootDir := fmt.Sprintf("Data/%s", videoMD5)
				os.MkdirAll(mRootDir, os.ModePerm)
				var buff bytes.Buffer
				buff.WriteString("Error:")
				buff.WriteString(err.Error())
				buff.WriteString("\r\n")
				buff.WriteString("VideoInfo:\r\n")
				if data, err := json.Marshal(&videoInfo); err == nil {
					buff.Write(data)
				}
				buff.WriteString("\r\n")
				ioutil.WriteFile(filepath.Join(mRootDir, "FAIL.log"), buff.Bytes(), os.ModePerm)
				api.Warn("[Task]Tasks that can't be completed:%s\r\nDetails:%s", videoMD5, err.Error())
			}
		} else {
			// 添加已完成记录
			p.mRecorder.AddRecord(videoInfo.ID, videoInfo)
		}
	}
}

func (p *QianjiaoVideoModule) workerExecTask(videoInfo *VideoInfo) error {
	defer api.ProtectError()
	var err error
	var resp *http.Response
	videoMD5 := videoInfo.GetMD5()

	// HTTP查询资源是否存在(去重)
	if api.IsTest() == false {
		if resp, err = http.Post(mConfigInfo.WebAPI+"video/query/exist", "text/plain", bytes.NewReader([]byte(videoInfo.VideoTitle))); err == nil {
			data, _ := ioutil.ReadAll(resp.Body)
			api.Info("[WebAPI]%s StateCode:%d Details:\n%s", videoMD5, resp.StatusCode, string(data))
			if string(data) == "exists" {
				return nil
			}
		} else {
			return fmt.Errorf("[WebAPI]%s Error:%s", videoMD5, err.Error())
		}
	}

	// 创建主目录
	mRootDir := fmt.Sprintf("Data/%s", videoMD5)
	os.MkdirAll(mRootDir, os.ModePerm)

	// 下载视频文件
	downloader := NewDownloaderM3U8(videoInfo.VideoURL, mRootDir, videoInfo.HeadersMap)
	if err = downloader.WaitDownload(); err != nil {
		return fmt.Errorf("[Download]%s Error:%s", videoMD5, err.Error())
	}

	// 测试时拦截
	if api.IsTest() {
		if detail, err := api.ConvertToMP4(filepath.Join(mRootDir, "output.m3u8"), filepath.Join(mRootDir, "output.mp4")); err != nil {
			api.Warn("ConvertToMP4 Error:%s Detail:\n", err.Error(), detail)
		} else {
			mSaveDir := fmt.Sprintf("Data/%s.mp4", videoMD5)
			os.MkdirAll(mSaveDir, os.ModePerm)
			os.Rename(filepath.Join(mRootDir, "output.mp4"), fmt.Sprintf("%s/output.mp4", mSaveDir))
			api.Info("[Clear]%s", mRootDir)
			os.RemoveAll(mRootDir)
		}
		return nil
	}

	// SCP上传到资源服务器
	var msg string
	if msg, err = api.SCP(mConfigInfo.SCPAddr, mConfigInfo.SCPPort, mConfigInfo.SCPUser, mConfigInfo.SCPPass, mRootDir, mConfigInfo.SavePath); err == nil {
		api.Info("[SCP]%s Done", videoMD5)
	} else {
		return fmt.Errorf("[SCP]%s Error:%s\ndetails:\n%s", videoMD5, err.Error(), msg)
	}

	// HTTP通知资源接口(入库)
	type UploadPostJSON struct {
		VideoName      string `json:"videoName"`
		MD5Name        string `json:"md5Name"`
		VideoType      string `json:"videoType"`
		VideoChildType string `json:"videoChildType"`
		Tags           string `json:"tags"`
	}
	var jsonRawData UploadPostJSON
	jsonRawData.VideoName = videoInfo.VideoTitle
	jsonRawData.MD5Name = videoMD5
	jsonRawData.VideoType = videoInfo.MajorCategory
	jsonRawData.VideoChildType = videoInfo.MinorCategory
	jsonRawData.Tags = videoInfo.Tags
	jsonByteData, _ := json.Marshal(&jsonRawData)
	if resp, err = http.Post(mConfigInfo.WebAPI+"video/upload", "application/json", bytes.NewReader(jsonByteData)); err == nil {
		data, _ := ioutil.ReadAll(resp.Body)
		api.Info("[WebAPI]%s StateCode:%d Details:\n%s", videoMD5, resp.StatusCode, string(data))
	} else {
		return fmt.Errorf("[WebAPI]%s Error:%s", videoMD5, err.Error())
	}

	// HTTP通知切片接口(处理缩略图以及封面)
	if resp, err = http.Get(fmt.Sprintf(mConfigInfo.WebAPI+"/video/videoHandler?videoPath=%s&md5Name=%s", url.QueryEscape(mConfigInfo.SavePath), videoMD5)); err == nil {
		data, _ := ioutil.ReadAll(resp.Body)
		result := string(data)
		api.Info("[WebAPI]%s videoHandler StateCode:%d Details:\n%s", videoMD5, resp.StatusCode, result)
		if result != "success" {
			return fmt.Errorf("[WebAPI]%s videoHandler result:%s", videoMD5, result)
		}
	} else {
		return fmt.Errorf("[WebAPI]%s videoHandler Error:%s", videoMD5, err.Error())
	}

	// 清空目录
	api.Info("[Clear]%s", mRootDir)
	os.RemoveAll(mRootDir)
	return nil
}

func init() {
	api.RegisterCapturer(new(QianjiaoVideoModule))
}
