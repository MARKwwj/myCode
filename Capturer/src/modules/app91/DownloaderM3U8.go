package app91

import (
	"TianlangCapturer/src/api"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// DownloaderM3U8 DownloaderM3U8
type DownloaderM3U8 struct {
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

// NewDownloaderM3U8 NewDownloaderM3U8
func NewDownloaderM3U8(url string, saveDir string, reqestHeader map[string]string) *DownloaderM3U8 {
	d := &DownloaderM3U8{
		mTaskURL:      url,
		mSavePath:     saveDir,
		mReqestHeader: reqestHeader,
	}
	return d
}

// WaitDownload WaitDownload
func (d *DownloaderM3U8) WaitDownload() error {

	// 下载第一个M3U8文件
	fileName := "output.m3u8"
	fileSavePath := filepath.Join(d.mSavePath, fileName)
	if d.downloadFile(d.mTaskURL, fileSavePath, nil) == false {
		return fmt.Errorf("download failed:%s", d.mTaskURL)
	}

	var mRootURL string
	if index := strings.LastIndex(d.mTaskURL, "/"); index != -1 {
		mRootURL = d.mTaskURL[:index+1]
	}

	// 解析文件
	mM3U8File := d.parseM3U8File(fileSavePath)
	if mM3U8File == nil || len(mM3U8File.mURLs) <= 0 {
		return nil
	}

	// 多版本M3U8（正常情况多分辨率文件后缀应为.m3u8，但是个别平台有时没加.m3u8后缀，所以这里以非 .ts 后缀做为多版本判断）
	if m3u8URL := mM3U8File.mURLs[0]; strings.LastIndex(m3u8URL, ".ts") == -1 {
		if m3u8URL[0] == '/' {
			if index := strings.Index(mRootURL[8:], "/"); index != -1 {
				mRootURL = mRootURL[:8+index+1]
			}
		}

		d.mTaskURL = mRootURL + m3u8URL
		if index := strings.LastIndex(d.mTaskURL, "/"); index != -1 {
			mRootURL = d.mTaskURL[:index+1]
		}
		if d.downloadFile(d.mTaskURL, fileSavePath, nil) == false {
			return fmt.Errorf("download failed:%s", d.mTaskURL)
		}

		// 解析文件
		mM3U8File = d.parseM3U8File(fileSavePath)
	}

	// 是否包含有效TS
	if mM3U8File == nil || len(mM3U8File.mURLs) <= 0 || strings.LastIndex(mM3U8File.mURLs[0], ".ts") < 0 {
		return nil
	}

	// 下载所有TS文件
	switch mM3U8File.mMethod {
	case "":
		// 更新M3U8
		m3u8Data, _ := ioutil.ReadFile(fileSavePath)
		m3u8DataArray := strings.Split(string(m3u8Data), "\n")
		m3u8File, _ := os.Create(fileSavePath)
		defer m3u8File.Close()
		for _, line := range m3u8DataArray {
			if len(line) > 0 && line[0] != '#' {
				line = line[strings.LastIndex(line, "/")+1:]
			}
			m3u8File.WriteString(line)
			m3u8File.WriteString("\n")

			// 计算总时间
			if index := strings.Index(line, "#EXTINF:"); index >= 0 {
				if index2 := strings.Index(line, ","); index2 >= 0 {
					strTime := line[index+8 : index2]
					if fTime, err := strconv.ParseFloat(strTime, 64); err == nil {
						d.mAllVideoFileTotalTime += fTime
					}
				}
			}
		}

		mTSFileRoot := mRootURL
		if tsURL := mM3U8File.mURLs[0]; tsURL[0] == '/' {
			if index := strings.Index(mRootURL[8:], "/"); index != -1 {
				mTSFileRoot = mRootURL[:8+index+1]
			}
		} else if tsURL := mM3U8File.mURLs[0]; (len(tsURL) > 7 && tsURL[:7] == "http://") || (len(tsURL) > 8 && tsURL[:8] == "https://") {
			mTSFileRoot = ""
		}

		for i, fileURL := range mM3U8File.mURLs {
			fileName := filepath.Base(fileURL)
			fileSavePath := filepath.Join(d.mSavePath, fileName)
			if d.downloadFile(mTSFileRoot+fileURL, fileSavePath, nil) == false {
				return fmt.Errorf("download failed:%s", mTSFileRoot+fileURL)
			}
			if api.IsTest() && i > 3 && false {
				break
			}
		}
	case "AES-128":

		mKeyFileRoot := mRootURL
		if mM3U8File.mKeyURL[0] == '/' {
			if index := strings.Index(mRootURL[8:], "/"); index != -1 {
				mKeyFileRoot = mRootURL[:8+index+1]
			}
		}

		// 加载密钥
		keySavePath := filepath.Join(d.mSavePath, "AES-128.txt")
		if d.downloadFile(mKeyFileRoot+mM3U8File.mKeyURL, keySavePath, nil) == false {
			return fmt.Errorf("download failed:%s", mKeyFileRoot+mM3U8File.mKeyURL)
		}
		keyData, _ := ioutil.ReadFile(keySavePath)
		os.Remove(keySavePath)

		// 更新M3U8
		m3u8Data, _ := ioutil.ReadFile(fileSavePath)
		m3u8DataArray := strings.Split(string(m3u8Data), "\n")
		m3u8File, _ := os.Create(fileSavePath)
		defer m3u8File.Close()
		filter := "#EXT-X-KEY:"
		filterLen := len(filter)
		for _, line := range m3u8DataArray {
			if line == "" || len(line) > filterLen && line[:filterLen] == filter {
				continue
			}
			if len(line) > 0 && line[0] != '#' {
				line = line[strings.LastIndex(line, "/")+1:]
			}
			m3u8File.WriteString(line)
			m3u8File.WriteString("\n")

			// 计算总时间
			if index := strings.Index(line, "#EXTINF:"); index >= 0 {
				if index2 := strings.Index(line, ","); index2 >= 0 {
					strTime := line[index+8 : index2]
					if fTime, err := strconv.ParseFloat(strTime, 64); err == nil {
						d.mAllVideoFileTotalTime += fTime
					}
				}
			}
		}

		mTSFileRoot := mRootURL
		if tsURL := mM3U8File.mURLs[0]; tsURL[0] == '/' {
			if index := strings.Index(mRootURL[8:], "/"); index != -1 {
				mTSFileRoot = mRootURL[:8+index+1]
			}
		} else if tsURL := mM3U8File.mURLs[0]; (len(tsURL) > 7 && tsURL[:7] == "http://") || (len(tsURL) > 8 && tsURL[:8] == "https://") {
			mTSFileRoot = ""
		}

		// 下载并解密TS文件
		block, _ := aes.NewCipher(keyData)
		blockmode := cipher.NewCBCDecrypter(block, keyData)
		nBlockSize := blockmode.BlockSize()
		mBlockData := make([]byte, nBlockSize)
		for i, fileURL := range mM3U8File.mURLs {
			fileName := filepath.Base(fileURL)
			fileSavePath := filepath.Join(d.mSavePath, fileName)
			if d.downloadFile(mTSFileRoot+fileURL, fileSavePath, func(tmpBuff *bytes.Buffer, r io.Reader, w io.Writer) (int, error) {
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
				return fmt.Errorf("download failed:%s", mTSFileRoot+fileURL)
			}
			if api.IsTest() && i > 3 && false {
				break
			}
		}
	default:
		api.Warn("Unknown methods: %s", mM3U8File.mMethod)
	}

	return nil
}

func (d *DownloaderM3U8) decodeAES128(inFilePath string, outFilePath string, key []byte) error {
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
	mMethod string
	mKeyURL string
	mURLs   []string
}

func (d *DownloaderM3U8) parseM3U8File(file string) *M3U8File {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	mM3U8File := &M3U8File{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 7 && line[:7] == "#EXT-X-" {
			temp := line[1:]
			var name, value string
			index := strings.Index(temp, ":")
			if index > 0 {
				name = temp[:index]
				value = temp[index+1:]
				if name == "EXT-X-KEY" {
					attrArray := strings.Split(value, ",")
					for _, attrKeyValue := range attrArray {
						keyValue := strings.Split(attrKeyValue, "=")
						switch keyValue[0] {
						case "METHOD":
							mM3U8File.mMethod = keyValue[1]
						case "URI":
							mM3U8File.mKeyURL = keyValue[1]
							if mM3U8File.mKeyURL[0] == '"' {
								mM3U8File.mKeyURL = mM3U8File.mKeyURL[1 : len(mM3U8File.mKeyURL)-1]
							}
						}
					}
				}
			}
		} else if len(line) > 0 && line[0] != '#' {
			mM3U8File.mURLs = append(mM3U8File.mURLs, line)
		}
		/* else if strings.LastIndex(line, ".ts") > 0 || strings.LastIndex(line, ".m3u8") > 0 {
			mM3U8File.mURLs = append(mM3U8File.mURLs, line)
		} */
	}
	return mM3U8File
}

// DataHandlerFunc DataHandlerFunc
type DataHandlerFunc func(buffer *bytes.Buffer, r io.Reader, w io.Writer) (int, error)

func (d *DownloaderM3U8) defaultDataHandler(buffer *bytes.Buffer, r io.Reader, w io.Writer) (int, error) {
	data := make([]byte, 1024)
	n, err := r.Read(data)
	if n > 0 {
		w.Write(data[:n])
	}
	return n, err
}

func (d *DownloaderM3U8) downloadFile(srcURL string, saveFilePath string, dataHandler DataHandlerFunc) bool {
	if api.IsDebug() || true {
		api.Debug("[DownloadTask]%s", srcURL)
	}

	for tryCount := 0; tryCount < 10; tryCount++ {
		if d.downloadFileWorker(srcURL, saveFilePath, dataHandler) {
			d.mAllVideoFileTotalSize += d.mTotalFileSize
			return true
		}
		api.Warn("[Download]%s Download failed, try again after 3 seconds", saveFilePath)
		time.Sleep(3 * time.Second)
	}
	return false
}

func (d *DownloaderM3U8) downloadFileWorker(srcURL string, saveFilePath string, dataHandler DataHandlerFunc) bool {
	api.Info("[Download]%s", saveFilePath)

	reqest, err := http.NewRequest("GET", srcURL, nil)
	for key, value := range d.mReqestHeader {
		reqest.Header.Add(key, value)
	}
	if d.mCurreDownloadURL == srcURL && d.mCurreFileSize > 0 {
		reqest.Header.Set("Range", fmt.Sprintf("bytes=%d-", d.mCurreFileSize))
	}
	if dataHandler == nil {
		dataHandler = d.defaultDataHandler
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

	if d.mCurreDownloadURL != srcURL {
		d.mCurreFileSize = 0
		d.mTotalFileSize = int64(resp.ContentLength)
		d.mCurreDownloadURL = srcURL
		os.MkdirAll(filepath.Dir(saveFilePath), os.ModePerm)
	}

	mSaveFile, err := os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		api.Error("[Error]%s:%s", saveFilePath, err.Error())
		return false
	}
	defer mSaveFile.Close()
	mSaveFile.Seek(d.mCurreFileSize, os.SEEK_SET)

	nCount := int64(0)
	var buffer bytes.Buffer
	for {
		n, err := dataHandler(&buffer, resp.Body, mSaveFile)
		if n > 0 {
			nCount += int64(n)
			d.mCurreFileSize += int64(n)
			d.mDownloadingByte += int64(n)
		}
		if err != nil {
			if err != io.EOF && api.IsDebug() {
				api.Warn("[Download]progress(%d/%d) fail \n%s", nCount, resp.ContentLength, err.Error())
			}
			break
		}
	}
	return d.mCurreFileSize >= d.mTotalFileSize
}
