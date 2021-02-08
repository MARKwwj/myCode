package apptianlang

import (
	"TianlangCapturer/src/api"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

// TianlangLiveModuleConfigInfo TianlangLiveModuleConfigInfo
type TianlangLiveModuleConfigInfo struct {
	WebAPI       string `json:"webAPI"`
	ListenPort   string `json:"listenPort"`
	PushInterval int    `json:"pushInterval"`
	DeviceUUID   string `json:"deviceUUID"`
}

// TianlangLiveModule TianlangLiveModule
type TianlangLiveModule struct {
	mDisguiser               *Disguiser
	mConfigInfo              TianlangLiveModuleConfigInfo
	mCategoryModuleMap       map[string]CategoryModuleInfo
	mCategoryModuleMapGuard  sync.Mutex
	mEnablePullModuleID      map[string]bool
	mEnablePullModuleIDGuard sync.Mutex
}

// Name Name
func (p *TianlangLiveModule) Name() string { return "Tianlang.Live" }

// PlatformCategoryModuleInfo PlatformCategoryModuleInfo
type PlatformCategoryModuleInfo struct {
	Platform string               `json:"platform"`
	Modules  []CategoryModuleInfo `json:"module"`
}

// CategoryModuleInfo CategoryModuleInfo
type CategoryModuleInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (p *TianlangLiveModule) loadConfig() bool {
	api.Info("[Config]ReloadUpdate->TianlangLiveModule")
	data, err := ioutil.ReadFile("TianlangLiveModule.cfg")
	if err == nil {
		json.Unmarshal(data, &p.mConfigInfo)
	} else {
		data, _ = json.MarshalIndent(&p.mConfigInfo, "", "\t")
		ioutil.WriteFile("TianlangLiveModule.cfg", data, os.ModePerm)
	}
	return true
}

// Run Run
func (p *TianlangLiveModule) Run() (err error) {
	//默认配置
	p.mConfigInfo.WebAPI = "http://192.168.100.105:8801"
	p.mConfigInfo.ListenPort = "9991"
	p.mConfigInfo.PushInterval = 60 * 3
	p.mConfigInfo.DeviceUUID = "AB54F0C52D81149E813394738A07E1F5305D8A00"
	p.mEnablePullModuleID = make(map[string]bool)
	p.loadConfig()

	//开启API监听
	go http.ListenAndServe("0.0.0.0:"+p.mConfigInfo.ListenPort, p)

	for {
		if err = p.permanentRun(); err != nil {
			api.Error(err.Error())
			time.Sleep(time.Minute)
			api.Info("Try permanentRun")
		}
	}
}

func (p *TianlangLiveModule) permanentRun() (err error) {
	//走登录流程
	p.mDisguiser = NewDisguiser(p.mConfigInfo.DeviceUUID)
	if err = p.mDisguiser.Login(); err != nil {
		api.Error("LoginError:%s", err.Error())
		return err
	}
	if err = p.mDisguiser.Version(); err != nil {
		api.Error("VersionError:%s", err.Error())
		return err
	}

	//请求所有类别模块信息，并推送数据到后台API
	if err = p.PushAllCategoryModule(); err != nil {
		api.Error("PushAllCategoryModule:%s", err.Error())
		return err
	}

	sleepTime := rand.Int31() % 30
	api.Warn("[Disguiser]Sleep(3h%dm)", sleepTime)
	time.Sleep(3*time.Hour + time.Duration(sleepTime)*time.Minute)
	return nil
}

// PushAllCategoryModule PushAllCategoryModule
func (p *TianlangLiveModule) PushAllCategoryModule() (err error) {
	//抓取所有类别模块
	var result *api.AnyValue
	mCategoryModuleMap := make(map[string]CategoryModuleInfo)
	for curPage := 1; true; curPage++ {
		for tryCount := 0; tryCount < 10; tryCount++ {
			if result, err = p.mDisguiser.PullLiveCategory(curPage, 30); err == nil {
				break
			}
			api.Error("PullLiveCategory:%s", err.Error())
			time.Sleep(time.Second)
		}
		if err != nil {
			return
		}
		liveList := result.Find("live_list").AsSlice()
		for _, value := range liveList {
			info := CategoryModuleInfo{
				ID:    value.Find("dz").AsString(),
				Title: value.Find("name").AsString(),
			}
			mCategoryModuleMap[info.ID] = info
			api.Info("ID:%s Title:%s", info.ID, info.Title)
		}
		if len(liveList) < 30 {
			break
		}
	}

	//保存类别模块数据
	mPlatformCategoryModuleInfo := PlatformCategoryModuleInfo{Platform: "Tianlang"}
	mPlatformCategoryModuleInfo.Modules = make([]CategoryModuleInfo, 0, len(mCategoryModuleMap))
	for _, value := range mCategoryModuleMap {
		mPlatformCategoryModuleInfo.Modules = append(mPlatformCategoryModuleInfo.Modules, value)
	}
	p.mCategoryModuleMapGuard.Lock()
	p.mCategoryModuleMap = mCategoryModuleMap
	p.mCategoryModuleMapGuard.Unlock()

	//推送给后台
	for tryCount := 0; tryCount < 3; tryCount++ {
		if err = p.notifyServer("POST", p.mConfigInfo.WebAPI+"/longVideo/live", &mPlatformCategoryModuleInfo); err == nil {
			break
		}
		api.Error("NotifyServer(longVideo/live) error:%s", err.Error())
		time.Sleep(time.Second * 3)
	}
	return nil
}

// RequestPullCategoryModule RequestPullCategoryModule
type RequestPullCategoryModule struct {
	Platform string   `json:"platform"`
	Modules  []string `json:"module"`
}

// ServeHTTP ServeHTTP
func (p *TianlangLiveModule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var result WebAPIResponse
	result.mResponseWriter = w

	if api.IsDebug() {
		api.Debug("[WebAPI-%s]%s%s", r.Method, r.Host, r.URL.String())
	}

	switch r.URL.Path {
	case "/api/pull":
		if err := p.WebAPIRequestPullCategoryModule(r, &result); err != nil {
			result.Write(500, err.Error())
			break
		}
	default:
		result.Write(404, "Invalid request")
	}
}

// WebAPIResponse WebAPIResponse
type WebAPIResponse struct {
	Code            int    `json:"code"`
	Message         string `json:"msg"`
	mResponseWriter http.ResponseWriter
}

// Write Write
func (w *WebAPIResponse) Write(code int, msg string) error {
	w.Code, w.Message = code, msg
	dataJSON, _ := json.Marshal(w)
	w.mResponseWriter.Write(dataJSON)
	return nil
}

// IsEnablePull IsEnablePull
func (p *TianlangLiveModule) IsEnablePull(id string) bool {
	p.mEnablePullModuleIDGuard.Lock()
	defer p.mEnablePullModuleIDGuard.Unlock()
	if enbale, ok := p.mEnablePullModuleID[id]; ok {
		return enbale
	}
	return false
}

// WebAPIRequestPullCategoryModule WebAPIRequestPullCategoryModule
func (p *TianlangLiveModule) WebAPIRequestPullCategoryModule(req *http.Request, res *WebAPIResponse) error {
	dataJSON, _ := ioutil.ReadAll(req.Body)
	if api.IsDebug() {
		api.Debug("[WebAPI]RequestPullCategoryModuleData:\n%s", hex.Dump(dataJSON))
	}

	var mRequestPullCategoryModule RequestPullCategoryModule
	if err := json.Unmarshal(dataJSON, &mRequestPullCategoryModule); err != nil {
		return err
	}

	// 取消所有
	p.mEnablePullModuleIDGuard.Lock()
	defer p.mEnablePullModuleIDGuard.Unlock()
	for key := range p.mEnablePullModuleID {
		p.mEnablePullModuleID[key] = false
	}

	// 按需启用
	p.mCategoryModuleMapGuard.Lock()
	defer p.mCategoryModuleMapGuard.Unlock()
	for _, id := range mRequestPullCategoryModule.Modules {
		if value, ok := p.mCategoryModuleMap[id]; ok {
			if p.IsEnablePull(id) == false {
				p.mEnablePullModuleID[id] = true
				go p.PushModuleLiveRoomList(value)
			}
			continue
		}
		api.Error("Invalid category module id：%s", id)
	}

	return res.Write(200, "Succeed")
}

// PlatformLiveRoomInfo PlatformLiveRoomInfo
type PlatformLiveRoomInfo struct {
	Platform     string         `json:"platform"`
	ModuleID     string         `json:"id"`
	LiveRoomList []LiveRoomInfo `json:"list"`
}

// LiveRoomInfo LiveRoomInfo
type LiveRoomInfo struct {
	Cover   string `json:"img"`
	Title   string `json:"title"`
	Address string `json:"address"`
}

// PushModuleLiveRoomList PushModuleLiveRoomList
func (p *TianlangLiveModule) PushModuleLiveRoomList(moduleInfo CategoryModuleInfo) {
	defer api.ProtectError()
	for p.IsEnablePull(moduleInfo.ID) {
		// api.Info("Timing pull(name:%s)", moduleInfo.ID)
		result, err := p.mDisguiser.PullLiveAnchor(moduleInfo.ID)
		if err != nil {
			api.Error("PullLiveAnchor(name:%s)->Error:%s", moduleInfo.ID, err.Error())
			time.Sleep(time.Minute)
			continue
		}

		liveAnchorList := result.Find("live_anchor_list").AsSlice()
		mLiveRoom := make(map[string]LiveRoomInfo, len(liveAnchorList))
		for _, anchor := range liveAnchorList {
			info := LiveRoomInfo{
				Cover:   anchor.Find("img").AsString(),
				Title:   anchor.Find("title").AsString(),
				Address: anchor.Find("address").AsString(),
			}
			mLiveRoom[info.Address] = info
		}

		var mPlatformLiveRoomInfo PlatformLiveRoomInfo
		mPlatformLiveRoomInfo.Platform = "Tianlang"
		mPlatformLiveRoomInfo.ModuleID = moduleInfo.ID
		for _, value := range mLiveRoom {
			mPlatformLiveRoomInfo.LiveRoomList = append(mPlatformLiveRoomInfo.LiveRoomList, value)
		}

		//推送给后台
		for tryCount := 0; tryCount < 3; tryCount++ {
			if err = p.notifyServer("POST", p.mConfigInfo.WebAPI+"/longVideo/LiveSeeJson", &mPlatformLiveRoomInfo); err == nil {
				break
			}
			api.Error("NotifyServer(longVideo/LiveSeeJson) error:%s", err.Error())
			time.Sleep(time.Second * 3)
		}

		time.Sleep(time.Second * time.Duration(p.mConfigInfo.PushInterval))
	}
}

func (p *TianlangLiveModule) notifyServer(method string, url string, data interface{}) (err error) {
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
	reqest.Header.Add("Content-Type", "application/json;charset=utf-8")

	//输出调试信息
	if api.IsDebug() {
		for key, value := range reqest.Header {
			api.Debug("[%s]:%v", key, value)
		}
		if data != nil {
			api.Debug("PostData:\n%s", hex.Dump(postData))
		}
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
	bodyData, _ := ioutil.ReadAll(resp.Body)

	if len(bodyData) <= 0 {
		api.Info("[Notify - RES-%s]%s Done", method, url)
	} else {
		api.Info("[Notify - RES-%s]%s Response results:\n%s", method, url, hex.Dump(bodyData))
	}

	return
}

func init() {
	api.RegisterCapturer(new(TianlangLiveModule))
}
