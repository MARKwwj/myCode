package appwangyue

import (
	"TianlangCapturer/src/api"
	"bytes"
	"compress/gzip"
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
	mUID             string
	mToken           string
	mHost            string
	mMobile          string
	mSaveAccountFunc func(uid, token string)
}

// NewDisguiser NewDisguiser
func NewDisguiser(uid, token string) *Disguiser {
	d := &Disguiser{mUID: uid, mToken: token}
	d.mHost = "http://s.huaji-b.com:57769"
	return d
}

// Login Login
func (d *Disguiser) Login() (err error) {
	api.Info("[Login]UID:%s Token:%s", d.mUID, d.mToken)

	var result Result
	if result, err = d.doServer("POST", "http://s.huaji-b.com:57769/api/public/?service=Home.getConfig", nil); err != nil {
		return
	}

	if result, err = d.doServer("GET", "http://s.huaji-b.com:57769/api/public/?service=Home.getDomainList&plat=0", nil); err != nil {
		return
	}
	if domainList := result.GetValue("data.info").AsSlice()[0].Find("data"); domainList != nil {
		d.mHost = strings.Split(domainList.AsString(), "\n")[0]
	}
	time.Sleep(100 * time.Millisecond)

	// 用现有的Token尝试登陆
	if len(d.mUID) > 0 && len(d.mToken) > 0 {
		if result, err = d.doServer("GET", "/api/public/?service=User.getBaseInfo&uid="+d.mUID+"&token="+d.mToken, nil); err != nil {
			return
		}
		if code := result.GetValue("data.code"); code != nil && code.AsInt() != 0 {
			// Token无效，需要重新登陆
			api.Info("Code:%d, Msg:%s", code.AsInt(), result.GetValue("data.msg").AsString())
		} else {
			return nil
		}
	}
	api.Info("令牌无效，需要账号登陆！")
	time.Sleep(300 * time.Millisecond)
	if len(d.mMobile) <= 0 {
		fmt.Print("请输入手机号以接收登陆验证码: ")
		var mobile string
		fmt.Scanln(&mobile)
		d.mMobile = mobile
	}

	// 请求发送手机验证码
	const SALT = "76576076c1f5f657b634e966c8836a06"
	url := fmt.Sprintf("/api/public/?service=Login.getCode&mobile=%s&sign=%s", d.mMobile, api.GetMD5("mobile="+d.mMobile+"&"+SALT))
	if result, err = d.doServer("GET", url, nil); err != nil {
		return
	}
	if code := result.GetValue("data.code"); code != nil && code.AsInt() != 0 {
		err = fmt.Errorf("Code:%d, Msg:%s", code.AsInt(), result.GetValue("data.msg").AsString())
		return
	}

	// 输入验证码:
	var passCode string
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("Please input code: ")
	fmt.Scanln(&passCode)
	url = fmt.Sprintf("/api/public/?service=Login.userLogin&user_login=%s&user_pass=%s&mobile_code=%s&source=Andriod:", d.mMobile, passCode, passCode)
	if result, err = d.doServer("GET", url, nil); err != nil {
		return
	}
	if code := result.GetValue("data.code"); code != nil && code.AsInt() != 0 {
		err = fmt.Errorf("Code:%d, Msg:%s", code.AsInt(), result.GetValue("data.msg").AsString())
		return
	}
	accountInfo := result.GetValue("data.info").AsSlice()[0]
	d.mUID = accountInfo.Find("id").AsString()
	d.mToken = accountInfo.Find("token").AsString()
	api.Info("[Logind]UID:%s Token:%s", d.mUID, d.mToken)
	if d.mSaveAccountFunc != nil {
		d.mSaveAccountFunc(d.mUID, d.mToken)
	}
	return nil
}

// GetHotList GetHotList
func (d *Disguiser) GetHotList(page int) (value *api.AnyValue, err error) {
	api.Info("[GetHot]UID:%s page:%d", d.mUID, page)

	var result Result
	sign := api.GetMD5("uid=" + d.mUID + "&token=" + d.mToken + "&yhyllsq888")
	url := fmt.Sprintf("/api/public/?service=Home.getHot&p=%d&sign=%s", page, sign)
	if result, err = d.doServer("GET", url, nil); err != nil {
		return
	}
	value = &(result.GetValue("data.info").AsSlice()[0])
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
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("User-Agent", "okhttp-okgo/jeasonlzy")

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

	if err = json.Unmarshal(bodyData, &result); err != nil {
		return
	}
	if result.Code != 200 {
		err = fmt.Errorf("Code:%d, Msg:%s", result.Code, result.Msg)
	}
	return
}
