package appwangyue

import (
	"TianlangCapturer/src/api"
	"sync"
	"time"
)

// WangyueLiveModule WangyueLiveModule
type WangyueLiveModule struct {
	mDisguiser    *Disguiser
	mAccountUID   string
	mAccountToken string
}

// Name Name
func (p *WangyueLiveModule) Name() string { return "Wangyue.Live" }

// Run Run
func (p *WangyueLiveModule) Run() (err error) {
	//初始配置
	if mConfig := api.GetLiveServiceConfig(); mConfig != nil {
		//默认配置
		mConfig.SetDefaultValue("WebAPIHost", "http://192.168.100.222:9999")
		mConfig.SetDefaultValue("PushInterval", 30)
		// mConfig.SetDefaultValue("AccountUID", "")
		// mConfig.SetDefaultValue("AccountToken", "")
		mConfig.SetDefaultValue("AccountUID", "10201324")
		mConfig.SetDefaultValue("AccountToken", "03fbede7442e4d9beba7e579315d37dd")
		//读取配置
		p.mAccountUID = mConfig.Find("AccountUID").AsString()
		p.mAccountToken = mConfig.Find("AccountToken").AsString()
	}

	//启用直播抓取服务
	return api.EnableLiveService("Wangyue", p.onGetPlatformLiveModule, p.onGetPlatformLiveRoomInfo)
}

func (p *WangyueLiveModule) onGetPlatformLiveModule() (result []*api.PlatformLiveModuleInfo, err error) {
	api.Info("[Pull]GetPlatformLiveModule(Platform:Wangyue)")
	defer api.Info("[Done]GetPlatformLiveModule(Platform:Wangyue)")

	//走登录流程
	p.mDisguiser = NewDisguiser(p.mAccountUID, p.mAccountToken)
	p.mDisguiser.mSaveAccountFunc = func(uid, token string) {
		if mConfig := api.GetLiveServiceConfig(); mConfig != nil {
			mConfig.SetValue("AccountUID", uid)
			mConfig.SetValue("AccountToken", token)
			mConfig.Flush()
		}
	}
	if err = p.mDisguiser.Login(); err != nil {
		api.Error("LoginError:%s", err.Error())
		return
	}
	time.Sleep(time.Second * 3)

	//该平台没有子模块分类，自填充一个子模块数据
	module := &api.PlatformLiveModuleInfo{Platform: "wangyue"}
	module.ModuleID = "wangyue"
	module.ModuleName = "望月"
	result = append(result, module)
	return
}

func (p *WangyueLiveModule) onGetPlatformLiveRoomInfo(module *api.PlatformLiveModuleInfo) (result []*api.PlatformLiveRoomInfo, err error) {
	api.Info("[Pull]GetPlatformLiveRoomInfo(Platform:%s, Module:%s)", module.Platform, module.ModuleID)
	defer api.Info("[Done]GetPlatformLiveRoomInfo(Platform:%s, Module:%s)", module.Platform, module.ModuleID)

	//拉取直播间信息
	var roomList []*api.PlatformLiveRoomInfo
	for page := 1; true; page++ {
		var resp *api.AnyValue
		if resp, err = p.mDisguiser.GetHotList(page); err != nil {
			return
		}
		list := resp.Find("list").AsSlice()
		for _, item := range list {
			roomInfo := &api.PlatformLiveRoomInfo{
				Title:   item.Find("title").AsString(),
				Cover:   item.Find("thumb").AsString(),
				Address: item.Find("pull").AsString(),
			}
			roomList = append(roomList, roomInfo)
		}
		if len(list) < 50 {
			break
		}
	}

	//筛选有效直播间
	var waitGroup sync.WaitGroup
	var resultGuard sync.Mutex
	for i := range roomList {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			info := roomList[index]
			if err := api.CheckLiveURLState(info.Address); err == nil {
				if api.IsDebug() {
					api.Warn("Add=>> Title:%s Address:%s", info.Title, info.Address)
				}
				resultGuard.Lock()
				result = append(result, info)
				resultGuard.Unlock()
			} else if api.IsDebug() {
				api.Warn("Invalid=>> Title:%s Address:%s Error:%s", info.Title, info.Address, err.Error())
			}
		}(i)
	}
	waitGroup.Wait()
	return
}

func init() {
	api.RegisterCapturer(new(WangyueLiveModule))
}
