package apptianlang

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
	"strings"
	"time"
)

// DownloaderM3U8 DownloaderM3U8
type DownloaderM3U8 struct {
	mTaskURL          string
	mSavePath         string
	mCurreFileSize    int64
	mTotalFileSize    int64
	mDownloadingByte  int64
	mCurreDownloadURL string
}

// NewDownloaderM3U8 NewDownloaderM3U8
func NewDownloaderM3U8(url string, saveDir string) *DownloaderM3U8 {
	d := &DownloaderM3U8{
		mTaskURL:  url,
		mSavePath: saveDir,
	}
	return d
}

// WaitDownload WaitDownload
func (d *DownloaderM3U8) WaitDownload() error {

	//下载第一个M3U8文件
	fileName := "output.m3u8"
	fileSavePath := filepath.Join(d.mSavePath, fileName)
	d.downloadFile(d.mTaskURL, fileSavePath, nil)

	var mHost string
	index := strings.Index(d.mTaskURL[8:], "/")
	mHost = d.mTaskURL[:8+index+1]

	//解析文件
	mM3U8File := d.parseM3U8File(fileSavePath)
	if mM3U8File == nil || len(mM3U8File.mURLs) <= 0 {
		return nil
	}

	//多版本M3U8
	if strings.LastIndex(mM3U8File.mURLs[0], ".m3u8") > 0 {
		d.downloadFile(mHost+mM3U8File.mURLs[0], fileSavePath, nil)
		//解析文件
		mM3U8File = d.parseM3U8File(fileSavePath)
	}

	//是否包含有效TS
	if mM3U8File == nil || len(mM3U8File.mURLs) <= 0 || strings.LastIndex(mM3U8File.mURLs[0], ".ts") < 0 {
		return nil
	}

	//下载所有TS文件
	switch mM3U8File.mMethod {
	case "":
		for _, fileURL := range mM3U8File.mURLs {
			fileName := filepath.Base(fileURL)
			fileSavePath := filepath.Join(d.mSavePath, fileName)
			d.downloadFile(mHost+fileURL, fileSavePath, nil)
		}
	case "AES-128":
		//加载密钥
		keySavePath := filepath.Join(d.mSavePath, "AES-128.txt")
		d.downloadFile(mHost+mM3U8File.mKeyURL, keySavePath, nil)
		keyData, _ := ioutil.ReadFile(keySavePath)
		os.Remove(keySavePath)

		//更新M3U8
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
			if line[0] != '#' {
				line = line[strings.LastIndex(line, "/")+1:]
			}
			m3u8File.WriteString(line)
			m3u8File.WriteString("\n")
		}

		//下载并解密TS文件
		block, _ := aes.NewCipher(keyData)
		blockmode := cipher.NewCBCDecrypter(block, keyData)
		nBlockSize := blockmode.BlockSize()
		mBlockData := make([]byte, nBlockSize)
		for _, fileURL := range mM3U8File.mURLs {
			fileName := filepath.Base(fileURL)
			fileSavePath := filepath.Join(d.mSavePath, fileName)
			d.downloadFile(mHost+fileURL, fileSavePath, func(tmpBuff *bytes.Buffer, r io.Reader, w io.Writer) (int, error) {
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
			})
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
		} else if strings.LastIndex(line, ".ts") > 0 || strings.LastIndex(line, ".m3u8") > 0 {
			mM3U8File.mURLs = append(mM3U8File.mURLs, line)
		}
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
	tryCount := 0
	for {
		if d.downloadFileWorker(srcURL, saveFilePath, dataHandler) {
			return true
		}
		if tryCount < 10 {
			break
		}
		tryCount++
		api.Debug("[Download]%s Download failed, try again after 3 seconds", saveFilePath)
		time.Sleep(3 * time.Second)
	}
	return false
}

func (d *DownloaderM3U8) downloadFileWorker(srcURL string, saveFilePath string, dataHandler DataHandlerFunc) bool {
	api.Debug("[Download]%s", saveFilePath)

	reqest, err := http.NewRequest("GET", srcURL, nil)
	if d.mCurreDownloadURL == srcURL && d.mCurreFileSize > 0 {
		reqest.Header.Set("Range", fmt.Sprintf("bytes=%d-", d.mCurreFileSize))
	}
	if dataHandler == nil {
		dataHandler = d.defaultDataHandler
	}

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
		os.MkdirAll(filepath.Dir(saveFilePath), os.ModePerm)
	}

	mSaveFile, err := os.Create(saveFilePath)
	if err != nil {
		api.Error("[Error]%s:%s", saveFilePath, err.Error())
		return false
	}
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