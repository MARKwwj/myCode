package app91

import (
	"TianlangCapturer/src/api"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Result Result
type Result struct {
	Data      string `json:"data"`
	Sign      string `json:"sign"`
	Code      int    `json:"errcode,omitempty"`
	Timestamp int    `json:"timestamp"`
	mValue    *api.AnyValue
}

// Find Find
func (r *Result) Find(key string) *api.AnyValue {
	if r.mValue == nil {
		mMap := make(map[string]interface{})
		json.Unmarshal([]byte(r.Data), &mMap)
		r.mValue = api.NewAnyValue(mMap)
	}
	return r.mValue.Find(key)
}

// Disguiser Disguiser
type Disguiser struct {
	mClient       *http.Client
	mAESKey       string
	mSignKey      string
	mLiveAESKey   string
	mLiveSignKey  string
	mUserToken    string
	mDevicesUUID  string
	mLiveSignData string
	mOauthID      string
}

// NewDisguiser NewDisguiser
func NewDisguiser(aesKey string, signKey string, liveKey string) *Disguiser {
	d := &Disguiser{
		mAESKey:      aesKey,
		mSignKey:     signKey,
		mLiveAESKey:  liveKey,
		mLiveSignKey: "kihfks3kjdhfksjh3kdjfs745dkslfh4",
	}
	return d
}

func (d *Disguiser) doServer(method string, url string, data map[string]interface{}) (result Result, err error) {
	//计算加密数据、时间戳
	type ReqParams struct {
		Data      string `json:"data"`
		Sign      string `json:"sign"`
		Timestamp string `json:"timestamp"`
	}
	var reqParams ReqParams
	data["app_type"] = "local"
	data["build_id"] = "a1000"
	data["oauth_type"] = "android"
	data["oauth_id"] = "b98d29603105ee965c333825d9fd7993"
	data["version"] = "3.9.1"
	data["token"] = ""
	postData, _ := json.Marshal(data)
	postData = AesCFBEncrypt(postData, []byte(d.mAESKey))
	reqParams.Data = fmt.Sprintf("%X", postData)
	reqParams.Timestamp = fmt.Sprintf("%010d", time.Now().Unix())

	//计算签名
	var buff bytes.Buffer
	buff.WriteString("data=")
	buff.WriteString(reqParams.Data)
	buff.WriteString("&timestamp=")
	buff.WriteString(reqParams.Timestamp)
	buff.WriteString(d.mSignKey)
	shaSum := fmt.Sprintf("%x", sha256.Sum256(buff.Bytes()))
	reqParams.Sign = fmt.Sprintf("%x", md5.Sum([]byte(shaSum)))

	//构造请求参数
	buff.Reset()
	buff.WriteString("timestamp=")
	buff.WriteString(reqParams.Timestamp)
	buff.WriteString("&data=")
	buff.WriteString(reqParams.Data)
	buff.WriteString("&sign=")
	buff.WriteString(reqParams.Sign)

	//构建请求头
	reqest, _ := http.NewRequest(method, url, bytes.NewReader(buff.Bytes()))
	reqest.Header.Add("Accept-Encoding", "gzip")
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Add("User-Agent", "okhttp-okgo/jeasonlzy")
	reqest.Header.Add("Connection", "Keep-Alive")
	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var logBuff bytes.Buffer
	if api.IsDebug() && false {
		logBuff.Reset()
		fmt.Fprintf(&logBuff, "[REQ-%s]%s\r\n", method, url)
		for key, value := range reqest.Header {
			fmt.Fprintf(&logBuff, "[%s]:%v\r\n", key, value)
		}
		fmt.Fprintf(&logBuff, "ContentData:\n%s", hex.Dump(buff.Bytes()))
		api.Debug(logBuff.String())
	}

	//发起请求
	mClient := &http.Client{}
	var resp *http.Response
	if resp, err = mClient.Do(reqest); err != nil {
		return
	}

	//读取应答数据
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

	//解密应答数据
	if len(bodyData) > 0 && resp.StatusCode == http.StatusOK {
		if err = json.Unmarshal(bodyData, &result); err != nil {
			api.Error(err.Error())
			return
		}
		data, _ := hex.DecodeString(result.Data)
		result.Data = string(AesCFBDecrypt(data, []byte(d.mAESKey)))
	}

	if api.IsDebug() && false {
		logBuff.Reset()
		fmt.Fprintf(&logBuff, "[RES-%s]%s\r\n", method, url)
		for key, value := range resp.Header {
			fmt.Fprintf(&logBuff, "[%s]:%v\r\n", key, value)
		}
		data, _ := json.Marshal(result)
		if len(data) > 256 {
			fmt.Fprintf(&logBuff, "ContentData:\n%s", hex.Dump(data[:256]))
		} else {
			fmt.Fprintf(&logBuff, "ContentData:\n%s", hex.Dump(data))
		}

		api.Debug(logBuff.String())
		ioutil.WriteFile("res.txt", logBuff.Bytes(), os.ModePerm)
	}
	return
}

func (d *Disguiser) doLiveServer(method string, url string, data map[string]interface{}) (result Result, err error) {
	//计算加密数据、时间戳
	type ReqParams struct {
		Data      string `json:"data"`
		Sign      string `json:"sign"`
		Timestamp string `json:"timestamp"`
	}
	var reqParams ReqParams
	data["oauth_type"] = "android"
	data["liveSignData"] = d.mLiveSignData
	data["oauth_id"] = d.mOauthID
	data["theme"] = "k91live"
	data["version"] = "1.1.1"
	data["token"] = ""
	postData, _ := json.Marshal(data)
	postData = AesCFBEncrypt(postData, []byte(d.mLiveAESKey))
	reqParams.Data = fmt.Sprintf("%X", postData)
	reqParams.Timestamp = fmt.Sprintf("%010d", time.Now().Unix())

	//计算签名
	var buff bytes.Buffer
	buff.WriteString("data=")
	buff.WriteString(reqParams.Data)
	buff.WriteString("&timestamp=")
	buff.WriteString(reqParams.Timestamp)
	buff.WriteString(d.mLiveSignKey)
	shaSum := fmt.Sprintf("%x", sha256.Sum256(buff.Bytes()))
	reqParams.Sign = fmt.Sprintf("%x", md5.Sum([]byte(shaSum)))

	//构造请求参数
	buff.Reset()
	buff.WriteString("timestamp=")
	buff.WriteString(reqParams.Timestamp)
	buff.WriteString("&data=")
	buff.WriteString(reqParams.Data)
	buff.WriteString("&sign=")
	buff.WriteString(reqParams.Sign)

	//构建请求头
	reqest, _ := http.NewRequest(method, url, bytes.NewReader(buff.Bytes()))
	reqest.Header.Add("Accept-Encoding", "gzip")
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	reqest.Header.Add("User-Agent", "okhttp-okgo/jeasonlzy")
	reqest.Header.Add("Connection", "Keep-Alive")
	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var logBuff bytes.Buffer
	if api.IsDebug() && false {
		logBuff.Reset()
		fmt.Fprintf(&logBuff, "[REQ-%s]%s\r\n", method, url)
		for key, value := range reqest.Header {
			fmt.Fprintf(&logBuff, "[%s]:%v\r\n", key, value)
		}
		fmt.Fprintf(&logBuff, "ContentData:\n%s", hex.Dump(buff.Bytes()))
		api.Debug(logBuff.String())
	}

	//发起请求
	mClient := &http.Client{}
	var resp *http.Response
	if resp, err = mClient.Do(reqest); err != nil {
		return
	}

	//读取应答数据
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

	//解密应答数据
	if len(bodyData) > 0 && resp.StatusCode == http.StatusOK {
		if err = json.Unmarshal(bodyData, &result); err != nil {
			api.Error(err.Error())
			return
		}
		data, _ := hex.DecodeString(result.Data)
		result.Data = string(AesCFBDecrypt(data, []byte(d.mLiveAESKey)))
	}

	if api.IsDebug() && false {
		logBuff.Reset()
		fmt.Fprintf(&logBuff, "[RES-%s]%s\r\n", method, url)
		for key, value := range resp.Header {
			fmt.Fprintf(&logBuff, "[%s]:%v\r\n", key, value)
		}
		data, _ := json.Marshal(result)
		if len(data) > 256 {
			fmt.Fprintf(&logBuff, "ContentData:\n%s", hex.Dump(data[:256]))
		} else {
			fmt.Fprintf(&logBuff, "ContentData:\n%s", hex.Dump(data))
		}

		api.Debug(logBuff.String())
		ioutil.WriteFile("res.txt", logBuff.Bytes(), os.ModePerm)
	}
	return
}
