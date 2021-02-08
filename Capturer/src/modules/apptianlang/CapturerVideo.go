package apptianlang

import (
	"TianlangCapturer/src/api"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

// ConfigInfo ConfigInfo
type ConfigInfo struct {
	WebAPI     string `json:"webAPI"`
	SCPAddr    string `json:"scpAddr"`
	SCPPort    string `json:"scpPort"`
	SCPUser    string `json:"scpUser"`
	SCPPass    string `json:"scpPass"`
	SavePath   string `json:"savePath"`
	MaxLimit   int    `json:"maxLimit"`
	DeviceUUID string `json:"deviceUUID"`
}

var (
	mConfigInfo ConfigInfo
)

// TianlangVideoModule TianlangVideoModule
type TianlangVideoModule struct {
}

// Name Name
func (p *TianlangVideoModule) Name() string { return "Tianlang.Video" }

// Run Run
func (p *TianlangVideoModule) Run() error {
	//默认配置
	mConfigInfo.WebAPI = "http://192.168.100.49:8805/"
	mConfigInfo.SCPAddr = "192.168.100.51"
	mConfigInfo.SCPPort = "22"
	mConfigInfo.SCPUser = "root"
	mConfigInfo.SCPPass = "Dd778899"
	mConfigInfo.SavePath = "~/scp-test"
	mConfigInfo.MaxLimit = 0
	mConfigInfo.DeviceUUID = "AB54F0C52D81149E813394738A07E1F5305D8A00"

	//开始捕获
	for InitDownloader(); p.loadConfig(); {
		disguiser := NewDisguiser(mConfigInfo.DeviceUUID)
		if err := disguiser.Login(); err != nil {
			api.Error("LoginError:%s", err.Error())
			return err
		}
		if err := disguiser.Version(); err != nil {
			api.Error("VersionError:%s", err.Error())
			return err
		}
		if err := disguiser.PullHome(); err != nil {
			api.Error("PullHomeError:%s", err.Error())
			return err
		}
		// for _, cateID := range disguiser.mVideoCategoryList {
		// 	if err := disguiser.PullAllVideoList(cateID); err != nil {
		// 		api.Error("PullVideoListError:%s", err.Error())
		// 	}
		// }

		for i := len(disguiser.mVideoCategoryList) - 1; i >= 0; i-- {
			cateID := disguiser.mVideoCategoryList[i]
			if err := disguiser.PullAllVideoList(cateID); err != nil {
				api.Error("PullVideoListError:%s", err.Error())
			}
		}

		if mConfigInfo.MaxLimit > 0 && mCountLimit >= mConfigInfo.MaxLimit {
			api.Warn("[Disguiser]Reach the limit(%d/%d), stop the capture", mCountLimit, mConfigInfo.MaxLimit)
			time.Sleep(time.Hour * 24 * 365)
		}

		api.Warn("[Disguiser]Sleep(1h30m)")
		time.Sleep(1*time.Hour + 30*time.Minute)
	}
	return nil
}

func (p *TianlangVideoModule) loadConfig() bool {
	api.Info("[Config]ReloadUpdate")
	data, err := ioutil.ReadFile("TianlangCapturer.cfg")
	if err == nil {
		json.Unmarshal(data, &mConfigInfo)
	} else {
		data, _ = json.MarshalIndent(&mConfigInfo, "", "\t")
		ioutil.WriteFile("TianlangCapturer.cfg", data, os.ModePerm)
	}
	return true
}

func init() {
	api.RegisterCapturer(new(TianlangVideoModule))
}
