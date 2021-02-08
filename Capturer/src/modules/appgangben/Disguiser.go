package appgangben

import (
	"TianlangCapturer/src/api"
	"bytes"
	"compress/gzip"
	"crypto/aes"
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
	mDeviceID  string
	mAES128Key string
}

// NewDisguiser NewDisguiser
func NewDisguiser(deviceID, aesKey string) *Disguiser {
	return &Disguiser{
		mAES128Key: aesKey,
		mDeviceID:  deviceID,
	}
}

// Login Login
func (x *Disguiser) Login() error {
	api.Info("[Disguiser]Login()")
	return nil
}

// CategoryInfo CategoryInfo
type CategoryInfo struct {
	ID   int
	Name string
	Icon string
}

// GetAllCategory GetAllCategory
func (x *Disguiser) GetAllCategory() (result []*CategoryInfo, err error) {
	api.Info("[Disguiser]GetAllCategory()")
	var reqest *api.AnyValue
	if reqest, err = x.doServer("GET", "https://10.zgxjdw.cn/api/category/all", nil); err != nil {
		return
	}

	if mCategoryItems := reqest.Find("data").AsSlice(); len(mCategoryItems) > 0 {
		result = make([]*CategoryInfo, len(mCategoryItems))
		for i, item := range mCategoryItems {
			info := new(CategoryInfo)
			info.ID = item.Find("id").AsInt()
			info.Name = item.Find("name").AsString()
			info.Icon = item.Find("icon").AsString()
			result[i] = info
		}
	}
	return
}

// VideoInfo VideoInfo
type VideoInfo struct {
	ID           string `json:"id"`
	Tags         string `json:"tags"`
	Title        string `json:"title"`
	Duration     int    `json:"duration"`
	CoverURL     string `json:"coverURL"`
	VideoURL     string `json:"videoURL"`
	Description  string `json:"description"`
	CategoryName string `json:"categoryName"`
}

// GetHotList GetHotList
func (x *Disguiser) GetHotList(categoryID, page, size int) (result []*VideoInfo, err error) {
	api.Info("[Disguiser]GetHotList(%d, %d, %d)", categoryID, page, size)
	var reqest *api.AnyValue
	if reqest, err = x.doServer("GET", "https://10.zgxjdw.cn/api/videos/popular/"+fmt.Sprintf("%d?page=%d&size=%d", categoryID, page, size), nil); err != nil {
		return
	}

	if mVideoItems := reqest.Find("data").AsSlice(); len(mVideoItems) > 0 {
		result = make([]*VideoInfo, len(mVideoItems))
		for i, item := range mVideoItems {
			info := new(VideoInfo)
			info.ID = fmt.Sprintf("Gangben.%d", item.Find("id").AsInt())
			info.Title = strings.TrimSpace(item.Find("title").AsString())
			info.Duration = item.Find("duration").AsInt()
			info.Description = item.Find("description").AsString()
			info.CoverURL = "https://new.guangfubao.com.cn/resource/image/" + item.Find("cover").AsString()

			if len(info.Description) <= 0 {
				info.Description = info.Title
			}

			mCode := item.Find("code").AsString()
			mDefinition := item.Find("definition").AsStringSlice()[0]
			info.VideoURL = "https://m4.hzmuqing.cn/resource/hls/" + mCode + "/" + mDefinition + "/" + mCode + ".m3u8"

			if tags := item.Find("tags").AsSlice(); len(tags) > 0 {
				tagNames := make([]string, 0, len(tags))
				for _, tag := range tags {
					tagNames = append(tagNames, tag.Find("name").AsString())
				}
				info.Tags = strings.Join(tagNames, ",")
			}

			result[i] = info
		}
	}
	return
}

func (x *Disguiser) doServer(method string, url string, data interface{}) (result *api.AnyValue, err error) {
	var dataBytes []byte
	var postData io.Reader
	if method == "POST" || method == "PUT" {
		d, _ := json.Marshal(data)
		dataBytes = d
		postData = bytes.NewReader(dataBytes)
	}

	reqest, _ := http.NewRequest(method, url, postData)
	reqest.Header.Set("Cookie", "X-Device-Id="+x.mDeviceID)
	reqest.Header.Set("Content-Type", "application/json")
	reqest.Header.Add("User-Agent", "okhttp-okgo/jeasonlzy")

	if api.IsDebug() {
		api.Debug("[REQ-%s]%s", method, url)
		for key, value := range reqest.Header {
			api.Debug("[%s]:%v", key, value)
		}
		if data != nil {
			api.Debug("PostData:\n%s", hex.Dump(dataBytes))
		}
	}

	var resp *http.Response
	mClient := &http.Client{Timeout: 10 * time.Second}
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

	if true {
		data := make([]byte, base64.StdEncoding.DecodedLen(len(bodyData)))
		n, err := base64.StdEncoding.Decode(data, bodyData)
		if err != nil {
			return nil, err
		}
		bodyData = data[:n]
		if api.IsDebug() {
			api.Debug("Base64.Decoded:\n%s", hex.Dump(bodyData))
		}
	}

	if true {
		block, _ := aes.NewCipher([]byte(x.mAES128Key))
		for i := len(bodyData)/block.BlockSize() - 1; i >= 0; i-- {
			index := block.BlockSize() * i
			block.Decrypt(bodyData[index:index+block.BlockSize()], bodyData[index:index+block.BlockSize()])
		}
		n := len(bodyData)
		unpadnum := int(bodyData[n-1])
		bodyData = bodyData[:n-unpadnum]
		if api.IsDebug() {
			api.Debug("AES-128.Decrypt:\n%s", hex.Dump(bodyData))
		}
	}

	resultMap := make(map[string]interface{})
	if err = json.Unmarshal(bodyData, &resultMap); err != nil {
		return
	}
	result = api.NewAnyValue(resultMap)

	if ecode := result.Find("code").AsInt(); ecode != 200 {
		err = fmt.Errorf("Code:%d, Message:%s", ecode, result.Find("message").AsString())
	}
	return
}
