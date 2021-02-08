package webmanhuawu

import (
	"TianlangCapturer/src/api"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// ConfigInfo ConfigInfo
type ConfigInfo struct {
	WebAPI    string `json:"webAPI"`
	SCPAddr   string `json:"scpAddr"`
	SCPPort   string `json:"scpPort"`
	SCPUser   string `json:"scpUser"`
	SCPPass   string `json:"scpPass"`
	SavePath  string `json:"savePath"`
	MachineID string `json:"machineId"`
}

var (
	mHistory    *api.Recorder
	mConfigInfo ConfigInfo
)

// ManhuawuModule ManhuawuModule
type ManhuawuModule struct {
}

// Name Name
func (p *ManhuawuModule) Name() string { return "Manhuawu" }

// Run Run
func (p *ManhuawuModule) Run() error {
	//默认配置
	mConfigInfo.WebAPI = "http://192.168.100.107:8900/"
	mConfigInfo.SCPAddr = "192.168.100.51"
	mConfigInfo.SCPPort = "22"
	mConfigInfo.SCPUser = "root"
	mConfigInfo.SCPPass = "Dd778899"
	mConfigInfo.SavePath = "/home/resources/wwwroot/limao/crypted/cartoonInfo"
	mConfigInfo.MachineID = "r1"

	//开始捕获
	var waitSync sync.WaitGroup
	for mHistory = api.GetHistory(p, p.onHistoryDataHandler); p.loadConfig(); p.sleepWait() {
		graphicNovel := GraphicNovel{mRootURL: "http://www.mhqwe.xyz"}
		for pageID := 1; true; pageID++ {
			bookList, err := graphicNovel.GetBookList(pageID)
			if err != nil {
				api.Error(err.Error())
				break
			}
			for _, book := range bookList {
				if _, ok := mHistory.FindRecord(book.GetMD5()); ok {
					continue
				}
				waitSync.Add(1)
				go func(book *BookInfo) {
					defer waitSync.Done()
					for tryCount := 0; tryCount < 5; tryCount++ {
						if graphicNovel.DownloadBookAllData(book) {
							mHistory.AddRecord(book.GetMD5(), book)
							return
						}
					}
					api.Error("Failed to download comic. Maximum retry has been reached. MD5:%s", book.GetMD5())
				}(book)
			}
			waitSync.Wait()
			if len(bookList) < 10 {
				break
			}
		}
	}
	return nil
}

func (p *ManhuawuModule) sleepWait() {
	api.Warn("[Disguiser]Sleep(1h30m)")
	time.Sleep(1*time.Hour + 30*time.Minute)
}

func (p *ManhuawuModule) onHistoryDataHandler(data []byte) (key, value interface{}, err error) {
	bookInfo := &BookInfo{}
	if err := json.Unmarshal(data, bookInfo); err != nil {
		return nil, nil, err
	}
	return bookInfo.GetMD5(), bookInfo, nil
}

func (p *ManhuawuModule) loadConfig() bool {
	api.Info("[Config]ReloadUpdate")
	data, err := ioutil.ReadFile("Webmanhuawu.cfg")
	if err == nil {
		json.Unmarshal(data, &mConfigInfo)
	} else {
		data, _ = json.MarshalIndent(&mConfigInfo, "", "\t")
		ioutil.WriteFile("Webmanhuawu.cfg", data, os.ModePerm)
	}
	return true
}

func init() {
	api.RegisterCapturer(new(ManhuawuModule))
}
