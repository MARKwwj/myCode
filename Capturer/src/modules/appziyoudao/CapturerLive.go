package appziyoudao

import (
	"TianlangCapturer/src/api"
	"sync"
	"time"
)

// ZiyoudaoLiveModule ZiyoudaoLiveModule
type ZiyoudaoLiveModule struct {
	mDisguiser   *Disguiser
	mDevicesUUID string
}

// Name Name
func (p *ZiyoudaoLiveModule) Name() string { return "Ziyoudao.Live" }

// Run Run
func (p *ZiyoudaoLiveModule) Run() (err error) {
	//初始配置
	if mConfig := api.GetLiveServiceConfig(); mConfig != nil {
		//默认配置
		mConfig.SetDefaultValue("WebAPIHost", "http://192.168.100.222:9999")
		mConfig.SetDefaultValue("PushInterval", 60*3)
		mConfig.SetDefaultValue("DevicesUUID", "283137DAB6AE4FC8ED91802DDB1D169F0AB3F5A5")
		//读取配置
		p.mDevicesUUID = mConfig.Find("DevicesUUID").AsString()
	}

	//启用直播抓取服务
	return api.EnableLiveService("Ziyoudao", p.onGetPlatformLiveModule, p.onGetPlatformLiveRoomInfo)
}

func (p *ZiyoudaoLiveModule) onGetPlatformLiveModule() (result []*api.PlatformLiveModuleInfo, err error) {
	api.Info("[Pull]GetPlatformLiveModule(Platform:Ziyoudao)")
	defer api.Info("[Done]GetPlatformLiveModule(Platform:Ziyoudao)")

	//走登录流程
	p.mDisguiser = NewDisguiser(p.mDevicesUUID)
	if err = p.mDisguiser.Login(); err != nil {
		api.Error("LoginError:%s", err.Error())
		return
	}
	time.Sleep(time.Second * 3)

	//拉取每页模块信息
	var liveInfos []LiveInfo
	for page := 1; true; page++ {

		if liveInfos, err = p.mDisguiser.GetLives(page, 30); err != nil {
			return
		}

		//填充结果
		for _, info := range liveInfos {
			if info.Status == 1 {
				module := &api.PlatformLiveModuleInfo{Platform: "Ziyoudao"}
				module.ModuleID = info.DZ
				module.ModuleName = info.Name
				result = append(result, module)
				// if api.IsDebug() {
				// 	api.Info("ID:%s Name:%s", module.ModuleID, module.ModuleName)
				// }
			}
		}

		if len(liveInfos) < 30 {
			break
		}

		time.Sleep(time.Second * 1)
	}
	return
}

func (p *ZiyoudaoLiveModule) onGetPlatformLiveRoomInfo(module *api.PlatformLiveModuleInfo) (result []*api.PlatformLiveRoomInfo, err error) {
	api.Info("[Pull]GetPlatformLiveRoomInfo(Platform:%s, Module:%s)", module.Platform, module.ModuleID)
	defer api.Info("[Done]GetPlatformLiveRoomInfo(Platform:%s, Module:%s)", module.Platform, module.ModuleID)

	//拉取直播间信息
	var roomList []LiveRoomInfo
	if roomList, err = p.mDisguiser.GetLiveRooms(module.ModuleID); err != nil {
		return
	}

	//筛选有效直播间
	var waitGroup sync.WaitGroup
	var resultGuard sync.Mutex
	for i := range roomList {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			info := &roomList[index]
			if err := api.CheckLiveURLState(info.Address); err == nil {
				if api.IsDebug() {
					api.Warn("Add=>> Title:%s Address:%s", info.Title, info.Address)
				}
				resultGuard.Lock()
				result = append(result, &api.PlatformLiveRoomInfo{
					Title:   info.Title,
					Cover:   info.Cover,
					Address: info.Address,
				})
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
	api.RegisterCapturer(new(ZiyoudaoLiveModule))
}
