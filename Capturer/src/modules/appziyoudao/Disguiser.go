package appziyoudao

import (
	"TianlangCapturer/src/api"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Disguiser Disguiser
type Disguiser struct {
	mHost        string
	mToken       string
	mDevicesUUID string
}

// NewDisguiser NewDisguiser
func NewDisguiser(devicesUUID string) *Disguiser {
	d := &Disguiser{mDevicesUUID: devicesUUID}
	d.mHost = "https://eksaweb.xyz:8888"
	return d
}

// Login Login
func (d *Disguiser) Login() (err error) {
	api.Info("[Login]DevicesUUID:%s", d.mDevicesUUID)

	type LoginArgs struct {
		Type        string `json:"type"`
		Time        string `json:"time"`
		SignName    string `json:"signName"`
		DevicesUUID string `json:"devices_uuid"`
	}
	var args LoginArgs
	args.Type = "android"
	args.Time = fmt.Sprintf("%d", time.Now().Unix())
	args.SignName = "111"
	args.DevicesUUID = d.mDevicesUUID

	var result Result
	if result, err = d.doServer("POST", "/api/v1/login", &args); err != nil {
		return
	}

	d.mToken = result.GetValue("data.token").AsString()
	api.Info("[Login]Token:%s", d.mToken)
	return
}

// LiveInfo LiveInfo
type LiveInfo struct {
	ID     int    `json:"id"`
	DZ     string `json:"dz"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

// GetLives GetLives
func (d *Disguiser) GetLives(page, limit int) (result []LiveInfo, err error) {
	api.Info("[GetLives]Page:%d Limit:%d", page, limit)

	type GetLiveRoomArgs struct {
		Page  string `json:"page"`
		Limit string `json:"limit"`
		Time  string `json:"time"`
	}
	var args GetLiveRoomArgs
	args.Page = fmt.Sprintf("%d", page)
	args.Limit = fmt.Sprintf("%d", limit)
	args.Time = fmt.Sprintf("%d", time.Now().Unix())

	var resp Result
	if resp, err = d.doServer("POST", "/api/v1/live", &args); err != nil {
		return
	}

	if liveList := resp.GetValue("data.live_list").AsSlice(); liveList != nil {
		result = make([]LiveInfo, len(liveList))
		for i, info := range liveList {
			result[i].ID = info.Find("id").AsInt()
			result[i].DZ = info.Find("dz").AsString()
			result[i].Name = info.Find("name").AsString()
			result[i].Status = info.Find("status").AsInt()
		}
	}
	return
}

// LiveRoomInfo LiveRoomInfo
type LiveRoomInfo struct {
	Title   string `json:"title"`
	Cover   string `json:"img"`
	Address string `json:"address"`
}

// GetLiveRooms GetLiveRooms
func (d *Disguiser) GetLiveRooms(name string) (result []LiveRoomInfo, err error) {
	api.Info("[GetLiveRooms]Name:%s", name)

	type GetLiveRoomsArgs struct {
		Name string `json:"name"`
		Time string `json:"time"`
	}
	var args GetLiveRoomsArgs
	args.Name = name
	args.Time = fmt.Sprintf("%d", time.Now().Unix())

	var resp Result
	if resp, err = d.doServer("POST", "/api/v1/live_anchor", &args); err != nil {
		return
	}

	if roomList := resp.GetValue("data.live_anchor_list").AsSlice(); roomList != nil {
		result = make([]LiveRoomInfo, len(roomList))
		for i, info := range roomList {
			result[i].Title = info.Find("title").AsString()
			result[i].Cover = info.Find("img").AsString()
			result[i].Address = info.Find("address").AsString()
		}
	}
	return
}

func (d *Disguiser) doServer(method string, url string, data interface{}) (result Result, err error) {
	if url[0] == '/' {
		url = d.mHost + url
	}
	var postData []byte
	var reqest *http.Request
	if data != nil {
		postData, _ = json.Marshal(data)
		reqest, _ = http.NewRequest(method, url, bytes.NewReader(postData))
	} else {
		postData = nil
		reqest, _ = http.NewRequest(method, url, nil)
		if strings.ToUpper(method) == "POST" {
			reqest.Header.Add("Content-Length", "0")
		}
	}
	reqest.Header.Add("Accept-Encoding", "gzip")
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Add("Content-Type", "application/json;charset=utf-8")
	reqest.Header.Add("User-Agent", "okhttp-okgo/jeasonlzy")
	if len(d.mToken) > 0 {
		reqest.Header.Add("Token", d.mToken)
	}

	if api.IsDebug() {
		api.Debug("[REQ-%s]%s", method, url)
		for key, value := range reqest.Header {
			api.Debug("[%s]:%v", key, value)
		}
		if data != nil {
			api.Debug("PostData:\n%s", hex.Dump(postData))
		}
	}

	mClient := &http.Client{Timeout: time.Second * 10, Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	var resp *http.Response
	if resp, err = mClient.Do(reqest); err != nil {
		return
	}

	bodyData, _ := ioutil.ReadAll(resp.Body)

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(bytes.NewReader(bodyData))
		if err != nil {
			panic(err)
		}
		var ungzipData bytes.Buffer
		_, err = io.Copy(&ungzipData, gzipReader)
		gzipReader.Close()
		if err != nil {
			panic(err)
		}
		bodyData = ungzipData.Bytes()
	}

	if api.IsDebug() {
		api.Debug(strings.Repeat("-", 64))
		api.Debug("[RES-%s]%s", method, reqest.URL)
		for key, value := range resp.Header {
			api.Debug("[%s]:%v", key, value)
		}
		if bodyData != nil {
			api.Debug("BodyData:\n%s", hex.Dump(bodyData))
		}
	}

	if bodyData != nil {
		bodyData, _ = base64.StdEncoding.DecodeString(string(bodyData))
		bodyData = api.AESDecryptCBC(bodyData, []byte("AbpDspjkVNnnXVwg"), []byte("tSKANBW7rhGW5hN3"))
	}

	if err = json.Unmarshal(bodyData, &result); err != nil {
		return
	}
	if result.Code != 200 {
		err = fmt.Errorf("Code:%d, Msg:%s", result.Code, result.Msg)
	}
	return
}
