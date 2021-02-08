package apptianlang

import (
	"TianlangCapturer/src/api"
	"fmt"
	"time"
)

// ReqVersion ReqVersion
type ReqVersion struct {
	Time string `json:"time"`
	Type string `json:"type"`
}

// Version Version
func (d *Disguiser) Version() (err error) {
	api.Info("[Version]")
	var req ReqVersion
	req.Type = "android"
	req.Time = fmt.Sprintf("%d", time.Now().Unix())

	var result Result
	if result, err = d.doServer("POST", "https://leiluohaoxiong.vip/api/v2/version", &req); err != nil {
		return
	}

	version := result.GetStringValue("info.version")
	api.Debug("LatestVersion:%s", version)
	return
}
