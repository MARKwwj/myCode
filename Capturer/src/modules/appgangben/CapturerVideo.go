package appgangben

import (
	"TianlangCapturer/src/api"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// CaptureGangbenVideoModuleConfig CaptureGangbenVideoModuleConfig
type CaptureGangbenVideoModuleConfig struct {
	WebAPI   string `json:"webAPI"`
	SCPAddr  string `json:"scpAddr"`
	SCPPort  string `json:"scpPort"`
	SCPUser  string `json:"scpUser"`
	SCPPass  string `json:"scpPass"`
	SavePath string `json:"savePath"`
}

// CaptureGangbenVideoModule CaptureGangbenVideoModule
type CaptureGangbenVideoModule struct {
	mAESKey string
	mConfig CaptureGangbenVideoModuleConfig
}

// Name Name
func (p *CaptureGangbenVideoModule) Name() string { return "Gangben.video" }

// Run Run
func (p *CaptureGangbenVideoModule) Run() error {
	// 读取配置信息
	api.LoadConfig(func() interface{} {
		p.mConfig.WebAPI = "http://192.168.100.113:8080/videos"
		p.mConfig.SCPAddr = ""
		p.mConfig.SCPPort = "22"
		p.mConfig.SCPUser = "root"
		p.mConfig.SCPPass = "123456"
		p.mConfig.SavePath = "~/longvideos"
		return &p.mConfig
	})

	// 初始化历史记录
	mRecorder := api.GetHistory(p, func(data []byte) (key, value interface{}, err error) {
		task := &VideoInfo{}
		if err := json.Unmarshal(data, task); err != nil {
			return nil, nil, err
		}
		return task.ID, task, nil
	})

	// 启动下载器
	api.EnableDownloader(6, func(taskPtr interface{}) {
		if task, ok := taskPtr.(*VideoInfo); ok {
			// 构建目录
			var err error
			mRootDir := fmt.Sprintf("./Data/%s/", task.ID)
			os.RemoveAll(mRootDir)
			os.MkdirAll(mRootDir, os.ModePerm)
			os.Mkdir(filepath.Join(mRootDir, "setcover"), os.ModeDir)
			mVideoKey := task.ID

			// 下载封面
			if _, err = api.DownloadFileToDisk(task.CoverURL, filepath.Join(mRootDir, "setcover", "cover.webp"), nil); err != nil {
				api.Error("[DownloadTask] ID:%s Download cover error:%s", task.ID, err.Error())
			}

			// 下载视频
			mDownloader := NewM3U8Downloader(task.VideoURL, mRootDir, nil)
			if err = mDownloader.WaitDownload(); err != nil {
				api.Error("[DownloadTask] ID:%s Download video error:%s", task.ID, err.Error())
				return
			}

			//HTTP通知资源接口(入库)
			type LongVideoJSON struct {
				Title         string `json:"video_title"`
				Intro         string `json:"video_intro"`
				Classify      string `json:"video_classify"`
				Duration      int    `json:"video_duration"`
				Tags          string `json:"video_tags"`
				VideoByteSize int    `json:"video_byte_size"`
			}
			var respData []byte
			if respData, err = p.doWebServer("POST", p.mConfig.WebAPI, &LongVideoJSON{
				Title:         task.Title,
				Intro:         task.Description,
				Classify:      "精选-" + task.CategoryName,
				Duration:      task.Duration,
				Tags:          task.Tags,
				VideoByteSize: int(mDownloader.mAllVideoFileTotalSize),
			}); err != nil {
				api.Error("Error post to longvideos: %s", err.Error())
				return
			}
			if len(respData) <= 0 {
				api.Info("Post to longvideos done.")
			} else {
				api.Info("Post to longvideos response: \n%s", hex.Dump(respData))
			}

			type LongVideoRespJSON struct {
				Code      int    `json:"code"`
				Message   string `json:"msg"`
				NotRepeat bool   `json:"not_repeat"`
				VideoID   int    `json:"video_id"`
			}
			var respJSON LongVideoRespJSON
			if err = json.Unmarshal(respData, &respJSON); err != nil {
				api.Error("Parse longvideo response data error: %s", err.Error())
				return
			}
			if respJSON.Code != 0 {
				api.Error("longvideo response.code == %d, msg: %s", respJSON.Code, respJSON.Message)
				return
			}

			// 检查视频是否已存在
			if respJSON.NotRepeat {
				task.ID = fmt.Sprintf("%d", respJSON.VideoID)
				mNewRootDir := fmt.Sprintf("./Data/%s/", task.ID)
				os.Rename(mRootDir, mNewRootDir)
				mRootDir = mNewRootDir

				// SCP上传到资源服务器
				if len(p.mConfig.SCPAddr) > 0 {
					if msg, err := api.SCP(p.mConfig.SCPAddr, p.mConfig.SCPPort, p.mConfig.SCPUser, p.mConfig.SCPPass, mRootDir, p.mConfig.SavePath); err == nil {
						api.Info("[SCP]%s Done", mVideoKey)
					} else {
						api.Error("[SCP]%s Error:%s\ndetails:\n%s", mVideoKey, err.Error(), msg)
						return
					}
				}
			}

			// 存档记录
			mRecorder.AddRecord(mVideoKey, task)

			// 清除数据
			if len(p.mConfig.SCPAddr) > 0 {
				api.Info("[Clear]%s", mVideoKey)
				os.RemoveAll(mRootDir)
			} else {
				api.Info("[Done]%s", mVideoKey)
			}
		}
	})

	// 开始捕获
	api.PermanentProtectCall(func() {
		defer time.Sleep(5 * time.Minute)
		defer api.Info("[Done]CaptureGangbenVideoModule.PermanentProtectCall()")

		disguiser := NewDisguiser("c18ba95901a6c376", "c38e60ea2c2994fd")
		if err := disguiser.Login(); err != nil {
			api.Error("[Disguiser]Login error:%s", err.Error())
			return
		}

		mCategoryItems, err := disguiser.GetAllCategory()
		if err != nil {
			api.Error("[Disguiser]GetAllCategory error:%s", err.Error())
			return
		}

		for _, category := range mCategoryItems {
			switch category.Name {
			case "高清无码":
				break
			default:
				continue
			}
			for page := 1; true; page++ {
				mVideoList, err := disguiser.GetHotList(category.ID, page, 6)
				if err != nil {
					api.Error("[Disguiser]GetHotList error:%s", err.Error())
					return
				}
				for i, item := range mVideoList {
					if _, ok := mRecorder.FindRecord(item.ID); !ok {
						item.CategoryName = category.Name
						api.PostDownloadTask(mVideoList[i])
					}
				}
				if len(mVideoList) < 6 {
					api.Info("[Done]CaptureGangbenVideoModule.GetHotList(%d)", category.ID)
					break
				}
				time.Sleep(time.Second)
			}
		}
	})
	return nil
}

func (p *CaptureGangbenVideoModule) doWebServer(method string, url string, data interface{}) (respData []byte, err error) {
	api.Info("[Notify - REQ-%s]%s", method, url)

	//准备好需要POST的数据
	var postData []byte
	if data != nil && (method == "POST" || method == "PUT") {
		switch v := data.(type) {
		case []byte:
			postData = v
		default:
			if postData, err = json.Marshal(v); err != nil {
				return
			}
		}
	}

	//设置好请求头
	reqest, _ := http.NewRequest(method, url, bytes.NewReader(postData))
	reqest.Header.Add("User-Agent", "okhttp-okgo/jeasonlzy")
	reqest.Header.Add("Content-Type", "application/json")

	//输出调试信息
	if api.IsDebug() {
		for key, value := range reqest.Header {
			api.Debug("[%s]:%v", key, value)
		}
		if data != nil {
			api.Debug("PostData:\n%s", hex.Dump(postData))
		}
	} else if postData != nil {
		data := postData
		// if len(postData) > 256 {
		// 	data = data[:256]
		// }
		api.Info("[Notify - REQ-%s]%s PostData:\n%s", method, url, hex.Dump(data))
	}

	//发起请求
	var resp *http.Response
	var mClient http.Client
	mClient.Timeout = time.Second * 30
	if resp, err = mClient.Do(reqest); err != nil {
		return
	}
	defer resp.Body.Close()

	//检查应答码
	if resp.StatusCode != 200 {
		api.Error("StatusCode:%d", resp.StatusCode)
		return
	}

	//读取应答数据
	if respData, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	return
}

func init() {
	api.RegisterCapturer(new(CaptureGangbenVideoModule))
}
