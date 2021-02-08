package appqianjiao

import (
	"TianlangCapturer/src/api"
	"bytes"
	"compress/gzip"
	"crypto/md5"
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
	mMechineCode   string
	mAuthorization string
}

// NewDisguiser NewDisguiser
func NewDisguiser(mechineCode string) *Disguiser {
	d := &Disguiser{
		mAuthorization: "1",
		mMechineCode:   mechineCode,
	}
	return d
}

// Login Login
func (d *Disguiser) Login() (err error) {
	type LoginArgs struct {
		Platform    int    `json:"platform"`
		MechineCode string `json:"mechineCode"`
	}
	var args LoginArgs
	args.Platform = 0
	args.MechineCode = d.mMechineCode

	var result *api.AnyValue
	if result, err = d.doServer("POST", "http://api.bgziyyj.cn/Base_Manage/UserInfo/UserLoginIng", &args); err != nil {
		return
	}
	d.mAuthorization = result.Find("Data").Find("Token").AsString()
	return
}

// TagInfo TagInfo
type TagInfo struct {
	ID       string `json:"Id"`
	VideoTag string `json:"VideoTag"`
	AreaID   string `json:"AreaId"`
	TagCover string `json:"TagCover"`
}

// GetTagByArea GetTagByArea
func (d *Disguiser) GetTagByArea(areaName string) (listInfo []TagInfo, err error) {
	type GetTagByAreaArgs struct {
		Area string `json:"Area"`
	}
	var args GetTagByAreaArgs
	args.Area = areaName

	var result *api.AnyValue
	if result, err = d.doServer("POST", "http://api.bgziyyj.cn/Base_Manage/TagVideo/GetTagByArea", &args); err != nil {
		return
	}
	datas := result.Find("Data").AsSlice()
	for _, data := range datas {
		var info TagInfo
		info.ID = data.Find("Id").AsString()
		info.AreaID = data.Find("AreaId").AsString()
		info.VideoTag = strings.TrimSpace(data.Find("VideoTag").AsString())
		info.TagCover = strings.TrimSpace(data.Find("TagCover").AsString())
		listInfo = append(listInfo, info)
	}
	return
}

// VideoInfo VideoInfo
type VideoInfo struct {
	ID            string            `json:"Id"`
	Tags          string            `json:"tags"`
	VideoTitle    string            `json:"VideoTitle"`
	VideoCover    string            `json:"VideoCover"`
	VideoURL      string            `json:"VideoUrl"`
	Headers       string            `json:"Headers"`
	MajorCategory string            `json:"MajorCategory"`
	MinorCategory string            `json:"MinorCategory"`
	RetryCount    int32             `json:"-"`
	HeadersMap    map[string]string `json:"-"`
}

// GetMD5 GetMD5
func (p *VideoInfo) GetMD5() string {
	value := fmt.Sprintf("Qianjiao.app|%s", p.ID)
	return fmt.Sprintf("%x", md5.Sum([]byte(value)))
}

// GetVideoByTag GetVideoByTag
func (d *Disguiser) GetVideoByTag(id string, page int) (listInfo []VideoInfo, err error) {
	type GetVideoByTagArgs struct {
		ID   string `json:"Id"`
		Next int    `json:"next"`
	}
	var args GetVideoByTagArgs
	args.ID = id
	args.Next = page

	var result *api.AnyValue
	if result, err = d.doServer("POST", "http://api.bgziyyj.cn/Base_Manage/VideoInfo/GetVideoByTag", &args); err != nil {
		return
	}
	datas := result.Find("Data").AsSlice()
	for _, data := range datas {
		var info VideoInfo
		info.ID = data.Find("Id").AsString()
		info.VideoTitle = strings.TrimSpace(data.Find("VideoTitle").AsString())
		info.VideoCover = strings.TrimSpace(data.Find("VideoCover").AsString())
		info.VideoURL = strings.TrimSpace(data.Find("VideoUrl").AsString())
		if headers := data.Find("Headers"); !headers.IsNull() {
			info.Headers = headers.AsString()
			info.HeadersMap = make(map[string]string)
			json.Unmarshal([]byte(info.Headers), &info.HeadersMap)
		}
		listInfo = append(listInfo, info)
	}
	return
}

func (d *Disguiser) doServer(method string, url string, data interface{}) (result *api.AnyValue, err error) {
	postData, _ := json.Marshal(data)
	reqest, _ := http.NewRequest(method, url, bytes.NewReader(postData))
	reqest.Header.Add("Accept-Encoding", "gzip")
	reqest.Header.Add("Authorization", "Bearer "+d.mAuthorization)
	reqest.Header.Add("Cache-Control", "no-cache")
	reqest.Header.Add("Connection", "Keep-Alive")
	reqest.Header.Add("User-Agent", "okhttp/3.12.1")
	reqest.Header.Add("Content-Type", "application/json; charset=utf-8")

	if api.IsDebug() {
		api.Debug("[REQ-%s]%s", method, url)
		for key, value := range reqest.Header {
			api.Debug("[%s]:%v", key, value)
		}
		if data != nil {
			api.Debug("PostData:\n%s", hex.Dump(postData))
		}
	}

	mClient := &http.Client{Timeout: 10 * time.Second}
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
		api.Debug("[RES-%s]%s", method, reqest.URL)
		for key, value := range resp.Header {
			api.Debug("[%s]:%v", key, value)
		}
		if bodyData != nil {
			api.Debug("BodyData:\n%s", hex.Dump(bodyData))
		}
	}

	resultMap := make(map[string]interface{})
	if err = json.Unmarshal(bodyData, &resultMap); err != nil {
		return
	}
	result = api.NewAnyValue(resultMap)

	if ecode := result.Find("ErrorCode").AsInt(); ecode != 0 {
		err = fmt.Errorf("ErrorCode:%d, Msg:%s", ecode, result.Find("Msg").AsString())
	}
	return
}

func (d *Disguiser) download(url string, headers string) (result []byte, err error) {
	reqest, _ := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Accept-Encoding", "gzip")
	reqest.Header.Add("Authorization", "Bearer "+d.mAuthorization)
	reqest.Header.Add("Cache-Control", "no-cache")
	reqest.Header.Add("Connection", "Keep-Alive")
	reqest.Header.Add("User-Agent", "okhttp/3.12.1")
	reqest.Header.Add("Content-Type", "application/json; charset=utf-8")
	if len(headers) > 0 {
		header := make(map[string]string)
		json.Unmarshal([]byte(headers), &header)
		for key, value := range header {
			reqest.Header.Add(key, value)
		}
	}

	mClient := &http.Client{Timeout: 10 * time.Second}
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
	result = bodyData
	return
}
