package app91

import (
	"TianlangCapturer/src/api"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Capture91VideoModuleConfig Capture91VideoModuleConfig
type Capture91VideoModuleConfig struct {
	WebAPI   string `json:"webAPI"`
	SCPAddr  string `json:"scpAddr"`
	SCPPort  string `json:"scpPort"`
	SCPUser  string `json:"scpUser"`
	SCPPass  string `json:"scpPass"`
	SavePath string `json:"savePath"`
	ResID    string `json:"resId"`
}

// Capture91VideoModule Capture91VideoModule
type Capture91VideoModule struct {
	mAESKey  string
	mSignKey string
	mLiveKey string
	mConfig  Capture91VideoModuleConfig
}

// Name Name
func (p *Capture91VideoModule) Name() string { return "91.video" }

// Run Run
func (p *Capture91VideoModule) Run() error {
	//默认配置
	p.mConfig.ResID = "r1"
	p.mConfig.WebAPI = "http://192.168.100.241:8000/smallvideos/"
	p.mConfig.SCPAddr = "192.168.100.100"
	p.mConfig.SCPPort = "22"
	p.mConfig.SCPUser = "root"
	p.mConfig.SCPPass = "123456"
	p.mConfig.SavePath = "~/smallvideos"

	//获取配置文件路径
	mConfigFile := os.Args[0]
	if index := strings.LastIndex(mConfigFile, "."); index > 0 {
		mConfigFile = mConfigFile[:index]
	}
	mConfigFile += ".cfg"
	//读取配置文件
	if data, err := ioutil.ReadFile(mConfigFile); err == nil {
		json.Unmarshal(data, &p.mConfig)
	}
	//回写配置文件
	if data, err := json.MarshalIndent(&p.mConfig, "", "\t"); err == nil {
		ioutil.WriteFile(mConfigFile, data, os.ModePerm)
	}

	type DownloadTask struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		PlayURL     string `json:"playURL"`
		ThumbImg    string `json:"thumbImg"`
		RemarkNames string `json:"remarkNames"`
	}

	mRecorder := api.GetHistory(p, func(data []byte) (key, value interface{}, err error) {
		task := &DownloadTask{}
		if err := json.Unmarshal(data, task); err != nil {
			return nil, nil, err
		}
		return fmt.Sprintf("91-%d", task.ID), task, nil
	})

	api.EnableDownloader(6, func(taskPtr interface{}) {
		if task, ok := taskPtr.(*DownloadTask); ok {
			// 构建目录
			mRootDir := fmt.Sprintf("./Data/%d/", task.ID)
			os.MkdirAll(mRootDir, os.ModePerm)
			os.Mkdir(filepath.Join(mRootDir, "setcover"), os.ModeDir)
			mVideoKey := fmt.Sprintf("91-%d", task.ID)

			// 下载封面
			if _, err := api.DownloadFileToDisk(task.ThumbImg, filepath.Join(mRootDir, "setcover", "cover.webp"), nil); err != nil {
				api.Error("[DownloadTask] ID:%d Download cover error:%s", task.ID, err.Error())
			}

			// 下载视频
			mDownloader := NewDownloaderM3U8(task.PlayURL, mRootDir, nil)
			if err := mDownloader.WaitDownload(); err != nil {
				api.Error("[DownloadTask] ID:%d Download video error:%s", task.ID, err.Error())
				return
			}

			//SCP上传到资源服务器
			if len(p.mConfig.SCPAddr) > 0 {
				if msg, err := api.SCP(p.mConfig.SCPAddr, p.mConfig.SCPPort, p.mConfig.SCPUser, p.mConfig.SCPPass, mRootDir, p.mConfig.SavePath); err == nil {
					api.Info("[SCP]%s Done", mVideoKey)
				} else {
					api.Error("[SCP]%s Error:%s\ndetails:\n%s", mVideoKey, err.Error(), msg)
					return
				}
			}

			//HTTP通知资源接口(入库)
			type SmallVideosJSON struct {
				Title       string `json:"title"`
				PicURL      string `json:"pic_url"`
				VideoURL    string `json:"video_url"`
				VideoTimes  int    `json:"video_times"` // 秒
				RemarkNames string `json:"remark_names"`
				VideoLength int    `json:"video_length"`
				ResID       string `json:"res_id"`
			}
			if respData, err := p.doWebServer("POST", p.mConfig.WebAPI, &SmallVideosJSON{
				Title:       task.Title,
				PicURL:      fmt.Sprintf("smallvideos/%d/setcover/cover.webp", task.ID),
				VideoURL:    fmt.Sprintf("smallvideos/%d/output.m3u8", task.ID),
				VideoTimes:  int(mDownloader.mAllVideoFileTotalTime),
				RemarkNames: task.RemarkNames,
				VideoLength: int(mDownloader.mAllVideoFileTotalSize),
				ResID:       p.mConfig.ResID,
			}); err != nil {
				api.Error("Error post to smallvideos: %s", err.Error())
			} else {
				if len(respData) <= 0 {
					api.Info("Post to smallvideos done.")
				} else {
					api.Info("Post to smallvideos response: \n%s", hex.Dump(respData))
				}
			}

			// 存档记录
			mRecorder.AddRecord(mVideoKey, task)

			// 清除数据
			api.Info("[Clear]%s", mVideoKey)
			// os.RemoveAll(mRootDir)
		}
	})

	api.PermanentProtectCall(func() {
		defer time.Sleep(10 * time.Second)
		disguiser := NewDisguiser(p.mAESKey, p.mSignKey, p.mLiveKey)

		//login-recommend
		testMap := make(map[string]interface{})
		testMap["mod"] = "system"
		testMap["code"] = "index"
		result, err := disguiser.doServer("POST", "http://interface.91apiapi.com:8080/api.php", testMap)
		if err != nil {
			api.Error("Error:%s", err.Error())
			return
		}

		// 收集数据
		os.Mkdir("Data", os.ModePerm)
		recommended := result.Find("data").Find("recommendedData").AsSlice()
		for _, value := range recommended {
			url := value.Find("playUrl").AsString()
			if index := strings.LastIndex(url, "&duration="); index > 0 {
				url = fmt.Sprintf("%s%d", url[:index+10], value.Find("duration").AsInt())
			} else {
				url += fmt.Sprintf("&duration=%d", value.Find("duration").AsInt())
			}
			// api.Debug("title:%d", len(value.Find("title").AsString()))
			// api.Debug("thumbImg: %s", value.Find("thumbImg").AsString())
			// api.Debug("playUrl: %s", url)
			// api.Debug("tags: %s", value.Find("tags").AsStringSlice())
			// api.Debug("-------------------------------------------------------------------------------")

			var remarkNames string
			if tags := value.Find("tags"); tags != nil {
				if tagSlice := tags.AsStringSlice(); len(tagSlice) > 0 {
					remarkNames = strings.Join(tagSlice, ",")
				}
			}

			id := value.Find("id").AsInt()
			if _, ok := mRecorder.FindRecord(fmt.Sprintf("91-%d", id)); !ok {
				api.PostDownloadTask(&DownloadTask{
					ID:          id,
					Title:       value.Find("title").AsString(),
					PlayURL:     url,
					ThumbImg:    value.Find("thumbImg").AsString(),
					RemarkNames: remarkNames,
				})
			}
		}
	})

	// api.Debug("len(Recommended) == %d  len(dataMap) == %d", len(recommended), len(dataMap))

	// //liveSignData
	// api.Debug("liveSignData")
	// testMap = make(map[string]interface{})
	// testMap["mod"] = "system"
	// testMap["code"] = "ping"
	// result, err = disguiser.doServer("POST", "http://interface.91apiapi.com:8080/api.php", testMap)
	// if err != nil {
	// 	return err
	// }
	// // api.Debug(result.Data)

	// //watchList
	// api.Debug("watchList")
	// testMap = make(map[string]interface{})
	// testMap["mod"] = "index"
	// testMap["code"] = "watchList"
	// testMap["items"] = "71698,"
	// testMap["timestamp"] = fmt.Sprintf("%010d", time.Now().Unix())
	// result, err = disguiser.doServer("POST", "http://interface.91apiapi.com:8080/api.php", testMap)
	// if err != nil {
	// 	return err
	// }
	// // api.Debug(result.Data)

	// //getNoticeAll
	// api.Debug("getNoticeAll")
	// testMap = make(map[string]interface{})
	// testMap["mod"] = "system"
	// testMap["code"] = "getNoticeAll"
	// result, err = disguiser.doServer("POST", "http://interface.91apiapi.com:8080/api.php", testMap)
	// if err != nil {
	// 	return err
	// }
	// // api.Debug(result.Data)

	// //user
	// api.Debug("user")
	// testMap = make(map[string]interface{})
	// testMap["mod"] = "user"
	// testMap["code"] = "index"
	// result, err = disguiser.doServer("POST", "http://interface.91apiapi.com:8080/api.php", testMap)
	// if err != nil {
	// 	return err
	// }
	// // api.Debug(result.Data)
	// disguiser.mLiveSignData = result.Find("liveSignData").AsString()
	// disguiser.mOauthID = result.Find("data").Find("info").Find("uuid").AsString()
	// api.Warn("CoreInfo\nOauthID:%s\nLiveSignData:%s", disguiser.mOauthID, disguiser.mLiveSignData)

	// //getConfig
	// api.Debug("getConfig")
	// testMap = make(map[string]interface{})
	// // testMap["oauth_type"] = "android"
	// // testMap["liveSignData"] = "0730C01E8D1F57B8517C8F6F565698D4ED9E8C75069B505F271970E5A6DC155B3200CFD9D9D1BACC12306C59ABEFB8EE1875AD9304C6E4DAEF0BB8387BC87AAB66F6131A1737B38200616F3E547D30AFAE6BF7D4CFAB2E7AD57D7EA92144AF53C7747E149C2E578FAAE478BBCA709CC8BAA3281DCE517E0CC241334D60DE2F65101C08080C4BEE5359662ED71CB583448C1CD8D2F051D4265D4A67DA443D9A884D9A1ED0693BF155F9479A53F63E810ECD1124537602E3CBE34BB32FA3"
	// // testMap["oauth_id"] = "e6508e4880fbcf061a0fdbe7bafce393"
	// // testMap["theme"] = "k91live"
	// // testMap["version"] = "1.1.1"
	// // testMap["token"] = ""
	// result, err = disguiser.doLiveServer("POST", "http://v2.ksapi002.me:2052/api.php/api/home/getConfig", testMap)
	// if err != nil {
	// 	return err
	// }
	// // api.Debug(result.Data)

	// //live
	// api.Debug("live")
	// testMap = make(map[string]interface{})
	// // testMap["oauth_type"] = "android"
	// // testMap["liveSignData"] = "0730C01E8D1F57B8517C8F6F565698D4ED9E8C75069B505F271970E5A6DC155B3200CFD9D9D1BACC12306C59ABEFB8EE1875AD9304C6E4DAEF0BB8387BC87AAB66F6131A1737B38200616F3E547D30AFAE6BF7D4CFAB2E7AD57D7EA92144AF53C7747E149C2E578FAAE478BBCA709CC8BAA3281DCE517E0CC241334D60DE2F65101C08080C4BEE5359662ED71CB583448C1CD8D2F051D4265D4A67DA443D9A884D9A1ED0693BF155F9479A53F63E810ECD1124537602E3CBE34BB32FA3"
	// // testMap["oauth_id"] = "e6508e4880fbcf061a0fdbe7bafce393"
	// // testMap["theme"] = "k91live"
	// // testMap["version"] = "1.1.1"
	// // testMap["token"] = ""
	// result, err = disguiser.doLiveServer("POST", "http://v2.ksapi002.me:2052/api.php/v2/live/index", testMap)
	// if err != nil {
	// 	return err
	// }
	// // api.Debug(result.Data)
	// liveList := result.Find("data").Find("live_list").AsSlice()
	// for _, live := range liveList {
	// 	api.Debug("UID:%d", live.Find("uid").AsInt())
	// 	api.Debug("Nickname:%d", len(live.Find("nickname").AsString()))
	// 	api.Debug("Title:%d", len(live.Find("title").AsString()))
	// 	api.Debug("Pull:%s", live.Find("pull").AsString())
	// 	api.Debug("-------------------------------------------------------------------------------")
	// }

	time.Sleep(time.Hour)
	return nil
}

func (p *Capture91VideoModule) doWebServer(method string, url string, data interface{}) (respData []byte, err error) {
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
	mCapture91Video := new(Capture91VideoModule)
	mCapture91Video.mAESKey = "e79465cfbb39cjdusimcuekd3b066a6e"
	mCapture91Video.mSignKey = "132f1537f85sjdpcm59f7e318b9epa51"
	mCapture91Video.mLiveKey = "ljhlksslgkjfhlksuo8472rju6p2od03"
	api.RegisterCapturer(mCapture91Video)
}

// AesCFBEncrypt AesCFBEncrypt
func AesCFBEncrypt(data []byte, key []byte) (encrypted []byte) {

	a2 := a(32, 16, nil, key, 0)
	key = a2[0]

	block, _ := aes.NewCipher(key)
	if len(data) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := a2[1]

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(data, data)

	var dataBuff bytes.Buffer
	dataBuff.Write(iv)
	dataBuff.Write(data)
	return dataBuff.Bytes()
}

// AesCFBDecrypt AesCFBDecrypt
func AesCFBDecrypt(encrypted []byte, key []byte) (decrypted []byte) {

	a2 := a(32, 16, nil, key, 0)
	key = a2[0]

	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

func a(i int, i2 int, bArr []byte, bArr2 []byte, i3 int) [][]byte {
	var (
		digest   []byte
		i4       int
		bArr3    = bArr
		bArr4    = bArr2
		i5       = i
		i6       = i2
		bArr8    []byte
		instance bytes.Buffer
		bArr5    = make([]byte, i5)
		bArr6    = make([]byte, i6)
		bArr7    = [][]byte{bArr5, bArr6}
	)

	if bArr4 == nil {
		return bArr7
	}
	i7 := i6
	i8 := 0
	i9 := 0
	i10 := i5
	i11 := 0
	for {
		instance.Reset()

		i12 := i11 + 1
		if i11 > 0 {
			instance.Write(bArr8)
		}

		instance.Write(bArr4)

		if bArr3 != nil {
			instance.Write(bArr3[0:8])
		}

		dv := md5.Sum(instance.Bytes())
		digest = dv[:]

		i13 := i3
		for i14 := 1; i14 < i13; i14++ {
			instance.Reset()
			instance.Write(digest)
			dv := md5.Sum(instance.Bytes())
			digest = dv[:]
		}
		if i10 > 0 {
			i4 = 0
			for i10 != 0 && i4 != len(digest) {
				bArr5[i8] = digest[i4]
				i10--
				i4++
				i8++
			}
		} else {
			i4 = 0
		}
		if i7 > 0 && i4 != len(digest) {
			for i7 != 0 && i4 != len(digest) {
				bArr6[i9] = digest[i4]
				i7--
				i4++
				i9++
			}
		}
		if i10 == 0 && i7 == 0 {
			break
		}
		i11 = i12
		bArr8 = digest
	}
	for i15 := 0; i15 < len(digest); i15++ {
		digest[i15] = 0
	}
	return bArr7
}
