package main

import (
	"TianlangCapturer/src/api"
	_ "TianlangCapturer/src/modules/app91"
	_ "TianlangCapturer/src/modules/appgangben"
	_ "TianlangCapturer/src/modules/appqianjiao"
	_ "TianlangCapturer/src/modules/apptianlang"
	_ "TianlangCapturer/src/modules/appwangyue"
	_ "TianlangCapturer/src/modules/appziyoudao"
	//_ "TianlangCapturer/src/modules/webhxsxs"
	_ "TianlangCapturer/src/modules/webmanhuawu"
	//_ "TianlangCapturer/src/modules/webyazhouseba"
)

var (
	// DefaultModule DefaultModule
	DefaultModule = "SoundNovel.Yazhouseba"
	// CurrentVersion CurrentVersion
	CurrentVersion = "v0.0.1"
)

func main() {
	api.StartAllCapturer(DefaultModule, CurrentVersion)
	// url := "rtmp://pull.wuagmd.cn/live/2961168_1606726076?txSecret=dbd6542655026d2474060156868f9e46&txTime=5FC4C56D"
	// if err := api.CheckLiveURLState(url); err != nil {
	// 	fmt.Printf("LiveState:%s", err.Error())
	// } else {
	// 	fmt.Printf("LiveState:%s", "ok")
	// }
}
