package apptianlang

import (
	"TianlangCapturer/src/api"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	mEnableDebug = true
)

// DownloadHandler DownloadHandler
type DownloadHandler interface {
	OnDownloadFail(downloader *Downloader)
	GetSaveFileName(downloader *Downloader) string
	OnDownloadComplete(downloader *Downloader)
	OnDownloadHandler(downloader *Downloader, data []byte) []byte
}

// TaskInfo TaskInfo
type TaskInfo struct {
	mURL     string
	mHandler DownloadHandler
}

// Downloader Downloader
type Downloader struct {
	mRootDir          string
	mQueue            chan TaskInfo
	mWaitSync         sync.WaitGroup
	mCurreFileSize    int64
	mTotalFileSize    int64
	mDownloadingByte  int64
	mCurreDownloadURL string
}

// NewDownloader NewDownloader
func NewDownloader(saveDir string) *Downloader {
	d := &Downloader{
		mRootDir: saveDir,
		mQueue:   make(chan TaskInfo, 16),
	}
	go d.run()
	return d
}

// AddDownload AddDownload
func (d *Downloader) AddDownload(url string, handler DownloadHandler) {
	d.mWaitSync.Add(1)
	d.mQueue <- TaskInfo{
		mURL:     url,
		mHandler: handler,
	}
}

// WaitDownload WaitDownload
func (d *Downloader) WaitDownload() {
	d.mWaitSync.Wait()
}

// run
func (d *Downloader) run() {
	for info := range d.mQueue {
		d.resetDownloadState()
		fileName := info.mHandler.GetSaveFileName(d)
		saveFilePath := filepath.Join(d.mRootDir, fileName)
		os.MkdirAll(filepath.Dir(saveFilePath), os.ModePerm)
		for i := 0; i < 30 && !d.downloadFile(info.mURL, saveFilePath, info.mHandler); i++ {
			if mEnableDebug {
				api.Debug("下载失败，1秒后重试：%s", info.mURL)
			}
			time.Sleep(time.Second * 1)
		}
	}
}

func (d *Downloader) resetDownloadState() {
	d.mCurreFileSize = 0
	d.mTotalFileSize = 0
	d.mDownloadingByte = 0
	d.mCurreDownloadURL = ""
}

func (d *Downloader) downloadFile(srcURL string, saveFilePath string, handler DownloadHandler) bool {
	reqest, err := http.NewRequest("GET", srcURL, nil)
	if d.mCurreDownloadURL == srcURL && d.mCurreFileSize > 0 {
		reqest.Header.Set("Range", fmt.Sprintf("bytes=%d-", d.mCurreFileSize))
	}

	// resp, err := http.Get(srcURL)
	mClient := &http.Client{}
	resp, err := mClient.Do(reqest)
	if err != nil {
		api.Error("[Error]%s", err.Error())
		return false
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		api.Error("[Error]StatusCode:%d", resp.StatusCode)
		return false
	}

	if d.mCurreDownloadURL != srcURL {
		d.mCurreFileSize = 0
		d.mTotalFileSize = int64(resp.ContentLength)
		d.mCurreDownloadURL = srcURL
	}

	mSaveFile, err := os.Create(saveFilePath)
	if err != nil {
		api.Error("[Error]%s:%s", saveFilePath, err.Error())
		return false
	}
	mSaveFile.Seek(d.mCurreFileSize, os.SEEK_SET)

	nCount := int64(0)
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			mSaveFile.Write(handler.OnDownloadHandler(d, buffer[:n]))
			nCount += int64(n)
			d.mCurreFileSize += int64(n)
			d.mDownloadingByte += int64(n)
		}
		if err != nil {
			if err != io.EOF && mEnableDebug {
				api.Warn("[Download]progress(%d/%d) fail \n%s", nCount, resp.ContentLength, err.Error())
			}
			break
		}
	}

	mSaveFile.Close()
	resp.Body.Close()
	return d.mCurreFileSize >= d.mTotalFileSize
}
