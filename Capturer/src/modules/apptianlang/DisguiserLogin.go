package apptianlang

import (
	"TianlangCapturer/src/api"
	"fmt"
	"time"
)

// ReqLogin ReqLogin
type ReqLogin struct {
	IP          string `json:"ip"`
	Time        string `json:"time"`
	Type        string `json:"type"`
	DevicesUUID string `json:"devices_uuid"`
}

// Login Login
func (d *Disguiser) Login() (err error) {
	api.Info("[Login]DevicesUUID:%s", d.mDevicesUUID)
	var req ReqLogin
	req.IP = "64.64.243.107"
	req.Type = "android"
	req.Time = fmt.Sprintf("%d", time.Now().Unix())
	req.DevicesUUID = d.mDevicesUUID

	var result Result
	if result, err = d.doServer("POST", "https://leiluohaoxiong.vip/api/v2/login", &req); err != nil {
		return
	}

	d.mToken = result.GetStringValue("token")
	d.mUserToken = result.GetStringValue("user_info.token")
	return nil
}
