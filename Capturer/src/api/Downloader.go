package api

import (
	"fmt"
	"sync"
	"time"
)

var (
	mDownloaderOnce      sync.Once
	mDownloadTaskHandler func(task interface{})
	mDownloadTaskQueue   chan interface{}
)

// EnableDownloader EnableDownloader
func EnableDownloader(num int, taskHandler func(task interface{})) {
	mDownloaderOnce.Do(func() {
		mDownloadTaskQueue = make(chan interface{}, num*6)
		mDownloadTaskHandler = taskHandler
		for i := 0; i < num; i++ {
			go func() {
				for task := range mDownloadTaskQueue {
					ProtectCall(func() {
						mDownloadTaskHandler(task)
					})
				}
			}()
		}
	})
}

// PostDownloadTask PostDownloadTask
func PostDownloadTask(task interface{}) {
	for {
		select {
		case mDownloadTaskQueue <- task:
			return
		default:
			Warn("[Downloader]PostDownloadTask fail: Queue is full!")
			time.Sleep(15 * time.Second)
		}
	}
}

// Downloader Downloader
type Downloader struct {
	mName string
}

// GetDownloader GetDownloader
func GetDownloader(module IModule) *Downloader {
	p := new(Downloader)
	p.mName = fmt.Sprintf("Downloader.%s", module.Name())
	return p
}

// AddTask AddTask
func (p *Downloader) AddTask(args ...interface{}) {
}

// type tagTask struct {
// 	mURL              string
// 	mSavePath         string
// 	mFailureCallback  func(response *http.Response, url string, savePath string)
// 	mCompleteCallback func(response *http.Response, url string, savePath string)
// }

// var (
// 	mDownloadTask = make(chan *tagTask, 4096)
// )

// func initDownloader() error {
// 	Info("Create download worker")
// 	for workerNum := runtime.NumCPU() * 2; workerNum >= 0; workerNum-- {
// 		go workerGoroutine()
// 	}

// 	Info("Load history record..")
// 	// if recordFile, err := os.Open("record.dat"); err == nil {
// 	// 	defer recordFile.Close()
// 	// 	reader := bufio.NewReader(recordFile)
// 	// 	for {
// 	// 		lineStr, err := reader.ReadString('\n')
// 	// 		if err != nil {
// 	// 			break
// 	// 		}
// 	// 		if data, err := base64.StdEncoding.DecodeString(lineStr); err == nil {
// 	// 			videoInfo := &model.VideoInfo{}
// 	// 			if model.Unmarshal(data, videoInfo) == nil && videoInfo.ID > 0 {
// 	// 				mKnownVideoInfo.Store(videoInfo.GetMD5(), videoInfo)
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }

// 	// mHistoryRecord = make(chan *model.VideoInfo, 128)
// 	// go func() {
// 	// 	var (
// 	// 		err  error
// 	// 		file *os.File
// 	// 	)
// 	// 	for videoInfo := range mHistoryRecord {
// 	// 		Info("[Doen]Task:%s", videoInfo.GetMD5())
// 	// 		data, _ := model.Marshal(videoInfo)
// 	// 		base64Data := base64.StdEncoding.EncodeToString(data)
// 	// 		if file, err = os.OpenFile("record.dat", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm); err == nil {
// 	// 			file.WriteString(base64Data)
// 	// 			file.WriteString("\r\n")
// 	// 			file.Close()
// 	// 		} else {
// 	// 			Error("[Record]Open 'record.dat' error:%s recordData:\n%s", err.Error(), base64Data)
// 	// 		}
// 	// 	}
// 	// }()
// 	return nil
// }

// func workerGoroutine() {
// 	for task := range mDownloadTask {
// 		retryCount := 10
// 		for retryCount > 0 && !workerExecTask(task) {
// 			retryCount--
// 		}
// 	}
// }

// func workerExecTask(task *tagTask) bool {
// 	defer ProtectError()

// 	return true
// }
