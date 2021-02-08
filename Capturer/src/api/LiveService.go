package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// 运行状态枚举值
const (
	StateStopd   = 0
	StateRuning  = 1
	StateStoping = 2
)

var (
	mOnce                    sync.Once
	mWebAPIHost              string
	mPushInterval            int
	mPlatformLiveModuleMap   = make(map[string]*PlatformLiveModule)
	mPlatformLiveModuleGuard sync.Mutex
)

// PlatformLiveModuleInfo PlatformLiveModuleInfo
type PlatformLiveModuleInfo struct {
	Platform   string `json:"platform"`
	ModuleID   string `json:"id"`
	ModuleName string `json:"name"`
}

// GetPlatformLiveModuleFunc GetPlatformLiveModuleFunc
type GetPlatformLiveModuleFunc func() ([]*PlatformLiveModuleInfo, error)

// PlatformLiveRoomInfo PlatformLiveRoomInfo
type PlatformLiveRoomInfo struct {
	Title   string `json:"title"`
	Cover   string `json:"img"`
	Address string `json:"address"`
	PlatID  int    `json:"platId"`
}

// GetPlatformLiveRoomInfoFunc GetPlatformLiveRoomInfoFunc
type GetPlatformLiveRoomInfoFunc func(module *PlatformLiveModuleInfo) ([]*PlatformLiveRoomInfo, error)

// EnableLiveService EnableLiveService
func EnableLiveService(platformName string, moduleCallback GetPlatformLiveModuleFunc, roomCallback GetPlatformLiveRoomInfoFunc) (err error) {

	//首次初始化
	mOnce.Do(initLiveService)

	//更新配置文件
	mLiveConfig.Flush()

	//检查平台模块是否已存在
	platformName = strings.ToLower(platformName)
	if _, ok := mPlatformLiveModuleMap[platformName]; ok {
		return fmt.Errorf("Duplicate platform name: %s", platformName)
	}

	//初始化平台模块数据
	mPlatformLiveModuleGuard.Lock()
	mPlatformLiveModule := &PlatformLiveModule{
		Platform:              platformName,
		mPullModuleState:      make(map[string]int8),
		mLiveRoomInfoCallback: roomCallback,
	}
	mPlatformLiveModuleMap[platformName] = mPlatformLiveModule
	mPlatformLiveModuleGuard.Unlock()

	//开启捕获逻辑
	for {
		if err = protectRun(mPlatformLiveModule, moduleCallback); err != nil {
			Error(err.Error())
			time.Sleep(time.Minute)
			Info("[LiveService]Platform:%s Try protectRun", platformName)
			continue
		}
		sleepTime := rand.Int31() % 30
		Warn("[LiveService]Platform:%s Sleep(3h%dm)", platformName, sleepTime)
		time.Sleep(3*time.Hour + time.Duration(sleepTime)*time.Minute)
	}
}

func protectRun(platform *PlatformLiveModule, callback GetPlatformLiveModuleFunc) (err error) {
	defer ProtectError()
	var moduleInfo []*PlatformLiveModuleInfo
	if moduleInfo, err = callback(); err != nil {
		return
	}
	platform.Modules = moduleInfo
	return
}

// PlatformLiveModule PlatformLiveModule
type PlatformLiveModule struct {
	Platform              string
	Modules               []*PlatformLiveModuleInfo
	mPullModuleState      map[string]int8
	mLiveRoomInfoCallback GetPlatformLiveRoomInfoFunc
}

func initLiveService() {

	//获取配置
	if mConfig := GetLiveServiceConfig(); mConfig != nil {
		mWebAPIHost = mConfig.Find("WebAPIHost").AsString()
		mPushInterval = mConfig.Find("PushInterval").AsInt()
	}

	//开启WebTest
	if IsDebug() {
		go http.ListenAndServe("0.0.0.0:9999", new(WebTest))
	}

	//开启处理
	go PermanentProtectCall(func() {
		defer func() {
			// 等待若干秒
			Info("[LiveService]Enter push interval(%ds) to wait.", mPushInterval)
			time.Sleep(time.Duration(mPushInterval) * time.Second)
			Info("[LiveService]Leave push interval to wait.")
		}()

		// 向后台查询需要拉取的平台模块
		var err error
		var respData []byte
		if respData, err = webServerAPI("GET", "/live/plat", nil); err != nil {
			Error("[LiveService]Get from %s/live/plat error:%s", mWebAPIHost, err.Error())
			return
		}
		if IsDebug() {
			Debug("[LiveService]Get from %s/live/plat data:\n%s", mWebAPIHost, hex.Dump(respData))
		}
		type PlatInfo struct {
			ID       int    `json:"platId"`
			Name     string `json:"platName"`
			Platform string `json:"-"`
			ModuleID string `json:"-"`
		}
		type ResponsePlat struct {
			Code    int        `json:"code"`
			Message string     `json:"msg"`
			Data    []PlatInfo `json:"data"`
		}
		var res ResponsePlat
		if err := json.Unmarshal(respData, &res); err != nil {
			Error("[LiveService]Unmarshal json data error: %s", err.Error())
			return
		}
		if res.Code != 0 {
			Error("[LiveService]Code != 0, Msg:%s", res.Message)
			return
		}

		// 调用各平台的拉取函数，并汇集所有结果
		var mAllLiveRoom []*PlatformLiveRoomInfo
		if platList := res.Data; len(platList) > 0 {
			for index := range platList {
				plat := &platList[index]
				args := strings.SplitN(plat.Name, "@", 2)
				plat.Platform = args[0]
				if len(args) >= 2 {
					plat.ModuleID = args[1]
				} else {
					plat.ModuleID = plat.Platform
				}
				plat.Platform = strings.ToLower(plat.Platform)
				plat.ModuleID = strings.ToLower(plat.ModuleID)

				mPlatformLiveModuleGuard.Lock()
				platform, _ := mPlatformLiveModuleMap[plat.Platform]
				mPlatformLiveModuleGuard.Unlock()

				if platform == nil {
					Error("[LiveService]Invalid Platform:%s", plat.Platform)
					continue
				}

				var moduleInfo *PlatformLiveModuleInfo
				for i, module := range platform.Modules {
					if module.ModuleID == plat.ModuleID {
						moduleInfo = platform.Modules[i]
						break
					}
				}
				if moduleInfo != nil {
					if liveRoomList, err := platform.mLiveRoomInfoCallback(moduleInfo); err == nil {
						for _, info := range liveRoomList {
							info.PlatID = plat.ID
							mAllLiveRoom = append(mAllLiveRoom, info)
						}
					} else {
						Error("[LiveService]GetLiveRoomInfo error! Platform:%s Module:%s Error:%s", plat.Platform, plat.ModuleID, err.Error())
					}
				} else {
					Warn("[LiveService]Invalid ModuleID! Platform:%s Module:%s", plat.Platform, plat.ModuleID)
				}
			}
		}

		// 推送给后台直播间总列表
		if len(mAllLiveRoom) > 0 {
			type LiveSeeJSONArgs struct {
				AllLiveRoom []*PlatformLiveRoomInfo `json:"list"`
			}
			args := LiveSeeJSONArgs{AllLiveRoom: mAllLiveRoom}
			if respData, err := webServerAPI("POST", "/live/LiveSeeJson", &args); err != nil {
				Error("[LiveService]Error pushing to server: %s", err.Error())
			} else {
				if len(respData) <= 0 {
					Info("[LiveService]Pushing to server done.")
				} else {
					Info("[LiveService]Pushing to server response: \n%s", hex.Dump(respData))
				}
			}
		} else {
			Warn("[LiveService]No live room")
		}
	})
}

// WebTest WebTest
type WebTest struct{}

// ServeHTTP ServeHTTP
func (p *WebTest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if IsDebug() {
		Debug("[WebTest-%s]%s%s", r.Method, r.Host, r.URL.String())
	}
	switch r.URL.Path {
	case "/live/plat":
		fmt.Fprint(w, `{"code": 0, "msg": "ok", "data": [ {"platId": 1, "platName": "Wangyue"} ]}`)
		break
	case "/live/LiveSeeJson":
		data, _ := ioutil.ReadAll(r.Body)
		Warn("收到推送的直播间数据:\n%s\n%s\n", hex.Dump(data), string(data))
		break
	}
}

func webServerAPI(method string, url string, data interface{}) (respData []byte, err error) {
	Info("[Notify - REQ-%s]%s", method, url)

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

	if url[0] == '/' {
		url = mWebAPIHost + url
	}

	//设置好请求头
	reqest, _ := http.NewRequest(method, url, bytes.NewReader(postData))
	reqest.Header.Add("User-Agent", "okhttp-okgo/jeasonlzy")
	reqest.Header.Add("Content-Type", "application/json;charset=utf-8")

	//输出调试信息
	if IsDebug() {
		for key, value := range reqest.Header {
			Debug("[%s]:%v", key, value)
		}
		if data != nil {
			Debug("PostData:\n%s", hex.Dump(postData))
		}
	} else if postData != nil {
		data := postData
		if len(postData) > 256 {
			data = data[:256]
		}
		Info("[Notify - REQ-%s]%s PostData:\n%s", method, url, hex.Dump(data))
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
		Error("StatusCode:%d", resp.StatusCode)
		return
	}

	//读取应答数据
	if respData, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	return
}

// LiveConfig LiveConfig
type LiveConfig struct {
	mValue *AnyValue
}

// SetValue SetValue
func (p *LiveConfig) SetValue(key string, value interface{}) {
	switch v := value.(type) {
	case int:
		value = float64(v)
	case int32:
		value = float64(v)
	case int64:
		value = float64(v)
	case uint:
		value = float64(v)
	case uint32:
		value = float64(v)
	case uint64:
		value = float64(v)
	}
	p.mValue.mMap[key] = value
}

// SetDefaultValue SetDefaultValue
func (p *LiveConfig) SetDefaultValue(key string, value interface{}) {
	if _, ok := p.mValue.mMap[key]; !ok {
		switch v := value.(type) {
		case int:
			p.mValue.mMap[key] = float64(v)
		case int32:
			p.mValue.mMap[key] = float64(v)
		case int64:
			p.mValue.mMap[key] = float64(v)
		case uint:
			p.mValue.mMap[key] = float64(v)
		case uint32:
			p.mValue.mMap[key] = float64(v)
		case uint64:
			p.mValue.mMap[key] = float64(v)
		default:
			p.mValue.mMap[key] = value
		}
	}
}

// Find Find
func (p *LiveConfig) Find(key string) *AnyValue {
	return p.mValue.Find(key)
}

// Flush Flush
func (p *LiveConfig) Flush() {
	mConfigFile := os.Args[0]
	if strings.LastIndex(mConfigFile, ".exe") > 0 {
		mConfigFile = mConfigFile[:len(mConfigFile)-4]
	}
	mConfigFile += ".cfg"
	if data, err := json.MarshalIndent(&p.mValue.mMap, "", "\t"); err == nil {
		ioutil.WriteFile(mConfigFile, data, os.ModePerm)
	}
}

var (
	mLiveConfig     LiveConfig
	mConfigInitOnce sync.Once
)

// GetLiveServiceConfig GetLiveServiceConfig
func GetLiveServiceConfig() *LiveConfig {
	mConfigInitOnce.Do(func() {
		//创建对象
		mLiveConfig.mValue = NewAnyValue(make(map[string]interface{}))

		//获取配置文件路径
		mConfigFile := os.Args[0]
		if strings.LastIndex(mConfigFile, ".exe") > 0 {
			mConfigFile = mConfigFile[:len(mConfigFile)-4]
		}
		mConfigFile += ".cfg"

		//读取配置文件
		if data, err := ioutil.ReadFile(mConfigFile); err == nil {
			json.Unmarshal(data, &mLiveConfig.mValue.mMap)
		}
	})
	return &mLiveConfig
}
