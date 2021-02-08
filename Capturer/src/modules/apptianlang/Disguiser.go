package apptianlang

import (
	"TianlangCapturer/src/api"
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

// Disguiser Disguiser
type Disguiser struct {
	mToken             string
	mUserToken         string
	mDevicesUUID       string
	mCategoryMap       sync.Map
	mVideoCategoryList []int32
}

// CategoryInfo CategoryInfo
type CategoryInfo struct {
	ID         int32
	Name       string
	SuperiorID int32
}

// NewDisguiser NewDisguiser
func NewDisguiser(devicesUUID string) *Disguiser {
	d := &Disguiser{
		// mClient:      &http.Client{},
		mDevicesUUID: devicesUUID,
	}
	return d
}

// SetCategoryName SetCategoryName
func (d *Disguiser) SetCategoryName(id int32, name string, superiorID int32) {
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, ",", "")
	name = strings.ReplaceAll(name, "|", "")

	info := CategoryInfo{
		ID:         id,
		Name:       name,
		SuperiorID: superiorID,
	}
	d.mCategoryMap.Store(id, info)
}

// GetCategoryName GetCategoryName
func (d *Disguiser) GetCategoryName(id int32) string {
	if value, ok := d.mCategoryMap.Load(id); ok {
		return (value.(CategoryInfo)).Name
	}
	return "unknown"
}

// GetCategorySuperiorName GetCategorySuperiorName
func (d *Disguiser) GetCategorySuperiorName(id int32) string {
	if value, ok := d.mCategoryMap.Load(id); ok {
		return d.GetCategoryName((value.(CategoryInfo)).SuperiorID)
	}
	return "unknown"
}

func (d *Disguiser) doServer(method string, url string, data interface{}) (result Result, err error) {
	postData, _ := json.Marshal(data)
	reqest, _ := http.NewRequest(method, url, bytes.NewReader(postData))
	reqest.Header.Add("User-Agent", "okhttp-okgo/jeasonlzy")
	reqest.Header.Add("Content-Type", "application/json;charset=utf-8")
	reqest.Header.Add("Accept-Encoding", "gzip")
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	if d.mToken != "" {
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

	mClient := &http.Client{}
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

	if data, err := base64.StdEncoding.DecodeString(string(bodyData)); err == nil {
		mAESKey, mAESIV := []byte("AbpDspjkVNnnXVwg"), []byte("tSKANBW7rhGW5hN3")
		var block cipher.Block
		block, _ = aes.NewCipher(mAESKey)
		blockmode := cipher.NewCBCDecrypter(block, mAESIV)
		blockmode.CryptBlocks(data, data)
		n := len(data)
		unpadnum := int(data[n-1])
		data = data[:n-unpadnum]
		bodyData = data
	}

	if api.IsDebug() {
		api.Debug("[RES-%s]%s%s", method, reqest.Host, reqest.URL)
		for key, value := range resp.Header {
			api.Debug("[%s]:%v", key, value)
		}
		if bodyData != nil {
			api.Debug("BodyData:\n%s", hex.Dump(bodyData))
		}
	}

	if err = json.Unmarshal(bodyData, &result); err != nil {
		return
	}
	if result.Code != 200 {
		err = fmt.Errorf("Code:%d, Msg:%s", result.Code, result.Msg)
	}
	return
}
