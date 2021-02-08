package apptianlang

import (
	"TianlangCapturer/src/api"
	"fmt"
	"strconv"
	"time"
)

// ReqLiveCategory ReqLiveCategory
type ReqLiveCategory struct {
	Page  string `json:"page"`
	Limit string `json:"limit"`
	Time  string `json:"time"`
}

// PullLiveCategory PullLiveCategory
func (d *Disguiser) PullLiveCategory(page, limit int) (value *api.AnyValue, err error) {
	api.Info("[PullLiveCategory]page:%d limit:%d", page, limit)
	var req ReqLiveCategory
	req.Page = strconv.Itoa(page)
	req.Limit = strconv.Itoa(limit)
	req.Time = fmt.Sprintf("%d", time.Now().Unix())

	var result Result
	if result, err = d.doServer("PUT", "https://leiluohaoxiong.vip/api/v2/live", &req); err != nil {
		return
	}
	if result.Code != 200 {
		err = fmt.Errorf("Code:%d Msg:%s", result.Code, result.Msg)
		return
	}

	value = api.NewAnyValue(result.Data)
	return
}

// ReqLiveAnchor ReqLiveAnchor
type ReqLiveAnchor struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

// PullLiveAnchor PullLiveAnchor
func (d *Disguiser) PullLiveAnchor(name string) (value *api.AnyValue, err error) {
	api.Info("[PullLiveAnchor]name:%s", name)
	var req ReqLiveAnchor
	req.Name = name
	req.Time = fmt.Sprintf("%d", time.Now().Unix())

	var result Result
	if result, err = d.doServer("POST", "https://leiluohaoxiong.vip/api/v2/live_anchor", &req); err != nil {
		return
	}
	if result.Code != 200 {
		err = fmt.Errorf("Code:%d Msg:%s", result.Code, result.Msg)
		return
	}

	value = api.NewAnyValue(result.Data)
	return
}
