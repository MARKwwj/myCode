package appgangben

import (
	"TianlangCapturer/src/api"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// M3U8Downloader M3U8Downloader
type M3U8Downloader struct {
	mTaskURL               string
	mSavePath              string
	mCurreFileSize         int64
	mTotalFileSize         int64
	mDownloadingByte       int64
	mCurreDownloadURL      string
	mReqestHeader          map[string]string
	mAllVideoFileTotalSize int64
	mAllVideoFileTotalTime float64
}

// NewM3U8Downloader NewM3U8Downloader
func NewM3U8Downloader(url string, saveDir string, reqestHeader map[string]string) *M3U8Downloader {
	x := &M3U8Downloader{
		mTaskURL:      url,
		mSavePath:     saveDir,
		mReqestHeader: reqestHeader,
	}
	return x
}

// WaitDownload WaitDownload
func (x *M3U8Downloader) WaitDownload() error {

	// 下载第一个M3U8文件
	fileName := "output.m3u8"
	fileSavePath := filepath.Join(x.mSavePath, fileName)
	if x.downloadFile(x.mTaskURL, fileSavePath, nil) == false {
		return fmt.Errorf("download failed:%s", x.mTaskURL)
	}

	var mRootURL string
	if index := strings.LastIndex(x.mTaskURL, "/"); index != -1 {
		mRootURL = x.mTaskURL[:index+1]
	}

	// 解析文件
	mM3U8File := x.parseM3U8File(fileSavePath)
	if mM3U8File == nil || mM3U8File.GetFileNum() <= 0 {
		return nil
	}

	// // 多版本M3U8（正常情况多分辨率文件后缀应为.m3u8，但是个别平台有时没加.m3u8后缀，所以这里以非 .ts 后缀做为多版本判断）
	// if m3u8URL := mM3U8File.mURLs[0]; strings.LastIndex(m3u8URL, ".ts") == -1 {
	// 	if m3u8URL[0] == '/' {
	// 		if index := strings.Index(mRootURL[8:], "/"); index != -1 {
	// 			mRootURL = mRootURL[:8+index+1]
	// 		}
	// 	}

	// 	x.mTaskURL = mRootURL + m3u8URL
	// 	if index := strings.LastIndex(x.mTaskURL, "/"); index != -1 {
	// 		mRootURL = x.mTaskURL[:index+1]
	// 	}
	// 	if x.downloadFile(x.mTaskURL, fileSavePath, nil) == false {
	// 		return fmt.Errorf("download failed:%s", x.mTaskURL)
	// 	}

	// 	// 解析文件
	// 	mM3U8File = x.parseM3U8File(fileSavePath)
	// }

	// 是否包含有效TS
	if mM3U8File == nil || mM3U8File.GetFileNum() <= 0 {
		return nil
	}

	// 更新M3U8
	m3u8File, _ := os.Create(fileSavePath)
	m3u8File.WriteString("#EXTM3U\n")
	for _, key := range mM3U8File.mIndexs {
		if key == "EXT-X-KEY" {
			continue
		}

		value, _ := mM3U8File.mFields[key]
		if len(value) == 1 {
			fmt.Fprintf(m3u8File, "#%s:%s\n", key, value[0])
		} else {
			for _, item := range value {
				if key == "EXTINF" {
					items := strings.Split(item, ",")
					items[1] = items[1][strings.LastIndex(items[1], "/")+1:]
					item = items[0] + ",\n" + items[1]
				}
				fmt.Fprintf(m3u8File, "#%s:%s\n", key, item)
			}
		}
	}
	m3u8File.WriteString("#EXT-X-ENDLIST\n")
	m3u8File.Close()

	var mCacheKeyURI string
	var mCacheKeyData []byte
	mNum := mM3U8File.GetFileNum()
	for i := 0; i < mNum; i++ {
		method, uri, iv, ts := mM3U8File.GetFileInfo(i)
		switch method {
		case "AES-128":
			mKeyFileRoot := mRootURL
			if uri[0] == '/' {
				if index := strings.Index(mRootURL[8:], "/"); index != -1 {
					mKeyFileRoot = mRootURL[:8+index+1]
				}
			}

			// 加载密钥
			if mCacheKeyURI != uri {
				mCacheKeyURI = uri
				keySavePath := filepath.Join(x.mSavePath, "AES-128.txt")
				if x.downloadFile(mKeyFileRoot+uri, keySavePath, nil) == false {
					return fmt.Errorf("download failed:%s", mKeyFileRoot+uri)
				}
				keyData, _ := ioutil.ReadFile(keySavePath)
				os.Remove(keySavePath)
				mCacheKeyData = keyData
			}

			mTSFileRoot := mRootURL
			if ts[0] == '/' {
				if index := strings.Index(mRootURL[8:], "/"); index != -1 {
					mTSFileRoot = mRootURL[:8+index+1]
				}
			} else if (len(ts) > 7 && ts[:7] == "http://") || (len(ts) > 8 && ts[:8] == "https://") {
				mTSFileRoot = ""
			}

			// 下载并解密TS文件
			block, _ := aes.NewCipher(mCacheKeyData)
			var ivData []byte
			if len(iv) > 0 {
				data, err := hex.DecodeString(iv[2:])
				if err != nil {
					panic(err)
				}
				ivData = data
			} else {
				ivData = mCacheKeyData
			}
			blockmode := cipher.NewCBCDecrypter(block, ivData)
			nBlockSize := blockmode.BlockSize()
			mBlockData := make([]byte, nBlockSize)
			fileName := filepath.Base(ts)
			fileSavePath := filepath.Join(x.mSavePath, fileName)
			if x.downloadFile(mTSFileRoot+ts, fileSavePath, func(tmpBuff *bytes.Buffer, r io.Reader, w io.Writer) (int, error) {
				n, err := r.Read(mBlockData)
				if n > 0 {
					tmpBuff.Write(mBlockData[:n])
					n = 0
					for tmpBuff.Len() >= nBlockSize {
						n, _ = tmpBuff.Read(mBlockData)
						blockmode.CryptBlocks(mBlockData, mBlockData)
						w.Write(mBlockData)
					}
				}
				return n, err
			}) == false {
				return fmt.Errorf("download failed:%s", mTSFileRoot+ts)
			}

			if api.IsTest() && i > 2 {
				return nil
			}

		default:
			panic(fmt.Sprintf("未处理的M3U8方法:%s", method))
		}
	}

	return nil
}

func (x *M3U8Downloader) decodeAES128(inFilePath string, outFilePath string, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	blockmode := cipher.NewCBCDecrypter(block, key)

	inFile, err := os.Open(inFilePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	nBlockSize := blockmode.BlockSize()
	mBlockData := make([]byte, nBlockSize)

	for {
		n, err := inFile.Read(mBlockData)
		if err != nil {
			break
		}
		block := mBlockData[:n]
		blockmode.CryptBlocks(block, block)
		outFile.Write(block)
	}

	// src = unpadding(src)
	return nil
}

func unpadding(src []byte) []byte {
	n := len(src)
	unpadnum := int(src[n-1])
	return src[:n-unpadnum]
}

// M3U8File M3U8File
type M3U8File struct {
	// mMethod string
	// mKeyURL string
	// mURLs   []string
	mIndexs []string
	mFields map[string][]string
}

// GetFileNum GetFileNum
func (x *M3U8File) GetFileNum() int {
	if value, ok := x.mFields["EXTINF"]; ok {
		return len(value)
	}
	return 0
}

// GetFileInfo GetFileInfo
func (x *M3U8File) GetFileInfo(index int) (method string, uri string, iv string, ts string) {
	if extKey, ok := x.mFields["EXT-X-KEY"]; ok {
		var value string
		if index >= len(extKey) {
			value = extKey[len(extKey)-1]
		} else {
			value = extKey[index]
		}
		values := strings.Split(value, ",")

		for _, item := range values {
			kv := strings.Split(item, "=")
			key := kv[0]
			value := kv[1]
			switch key {
			case "METHOD":
				method = value
			case "URI":
				uri = value
				if n := len(uri); n > 0 && uri[0] == '"' && uri[n-1] == '"' {
					uri = uri[1 : n-1]
				}
			case "IV":
				iv = value
			default:
				panic(fmt.Sprintf("未知Key：%s", key))
			}
		}
	}

	if extInf, ok := x.mFields["EXTINF"]; ok {
		ts = strings.Split(extInf[index], ",")[1]
		ts = strings.TrimSpace(ts)
	}
	return
}

func (x *M3U8Downloader) parseM3U8File(file string) *M3U8File {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	mM3U8File := &M3U8File{mFields: make(map[string][]string)}
	lines := strings.Split(string(data), "\n")
	if len(lines) <= 0 || strings.TrimSpace(lines[0]) != "#EXTM3U" {
		return nil
	}

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if len(line) > 4 && line[:4] == "#EXT" {
			temp := line[1:]
			if index := strings.Index(temp, ":"); index > 0 {
				name := temp[:index]
				value := temp[index+1:]
				if name == "EXTINF" {
					value += strings.TrimSpace(lines[i+1])
				}

				needAdd := true
				for _, item := range mM3U8File.mIndexs {
					if item == name {
						needAdd = false
						break
					}
				}
				if needAdd {
					mM3U8File.mIndexs = append(mM3U8File.mIndexs, name)
				}

				mM3U8File.mFields[name] = append(mM3U8File.mFields[name], value)
				if api.IsDebug() {
					api.Debug("%s:%s", name, value)
				}
			}
		}
	}
	return mM3U8File
}

// DataHandlerFunc DataHandlerFunc
type DataHandlerFunc func(buffer *bytes.Buffer, r io.Reader, w io.Writer) (int, error)

func (x *M3U8Downloader) defaultDataHandler(buffer *bytes.Buffer, r io.Reader, w io.Writer) (int, error) {
	data := make([]byte, 1024)
	n, err := r.Read(data)
	if n > 0 {
		w.Write(data[:n])
	}
	return n, err
}

func (x *M3U8Downloader) downloadFile(srcURL string, saveFilePath string, dataHandler DataHandlerFunc) bool {
	api.Debug("[DownloadTask]%s", srcURL)

	for tryCount := 0; tryCount < 10; tryCount++ {
		if x.downloadFileWorker(srcURL, saveFilePath, dataHandler) {
			x.mAllVideoFileTotalSize += x.mTotalFileSize
			return true
		}
		api.Warn("[Download]%s Download failed, try again after 3 seconds", saveFilePath)
		time.Sleep(3 * time.Second)
	}
	return false
}

func (x *M3U8Downloader) downloadFileWorker(srcURL string, saveFilePath string, dataHandler DataHandlerFunc) bool {
	api.Info("[Download]%s", saveFilePath)

	reqest, err := http.NewRequest("GET", srcURL, nil)
	for key, value := range x.mReqestHeader {
		reqest.Header.Add(key, value)
	}
	if x.mCurreDownloadURL == srcURL && x.mCurreFileSize > 0 {
		reqest.Header.Set("Range", fmt.Sprintf("bytes=%d-", x.mCurreFileSize))
	}
	if dataHandler == nil {
		dataHandler = x.defaultDataHandler
	}

	mClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := mClient.Do(reqest)
	if err != nil {
		api.Error("[Error]%s", err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		api.Error("[Error]StatusCode:%d", resp.StatusCode)
		return false
	}

	if x.mCurreDownloadURL != srcURL {
		x.mCurreFileSize = 0
		x.mTotalFileSize = int64(resp.ContentLength)
		x.mCurreDownloadURL = srcURL
		os.MkdirAll(filepath.Dir(saveFilePath), os.ModePerm)
	}

	mSaveFile, err := os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		api.Error("[Error]%s:%s", saveFilePath, err.Error())
		return false
	}
	defer mSaveFile.Close()
	mSaveFile.Seek(x.mCurreFileSize, os.SEEK_SET)

	nCount := int64(0)
	var buffer bytes.Buffer
	for {
		n, err := dataHandler(&buffer, resp.Body, mSaveFile)
		if n > 0 {
			nCount += int64(n)
			x.mCurreFileSize += int64(n)
			x.mDownloadingByte += int64(n)
		}
		if err != nil {
			if err != io.EOF && api.IsDebug() {
				api.Warn("[Download]progress(%d/%d) fail \n%s", nCount, resp.ContentLength, err.Error())
			}
			break
		}
	}
	return x.mCurreFileSize >= x.mTotalFileSize
}
