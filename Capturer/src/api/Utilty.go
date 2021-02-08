package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// IsTsSuffixFile IsTsSuffixFile
func IsTsSuffixFile(filepath string) bool {
	n := len(filepath)
	return n > 3 && filepath[n-3:] == ".ts"
}

// IsValidTSFile IsValidTSFile
func IsValidTSFile(file string) bool {
	in, err := os.Open(file)
	defer in.Close()
	if err == nil {
		data := make([]byte, 8)
		if in.Read(data); data[0] == 0x47 {
			return true
		}
	}
	return false
}

// ParseM3u8File ParseM3u8File
func ParseM3u8File(file string) []string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	var tsList []string
	lines := strings.Split(string(data), "\n")
	for _, name := range lines {
		if len(name) > 3 && name[len(name)-3:] == ".ts" {
			tsList = append(tsList, name)
		}
	}
	return tsList
}

// CalculateTsXORKey CalculateTsXORKey
func CalculateTsXORKey(file string) []byte {
	in, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer in.Close()
	data := make([]byte, 8)
	if in.Read(data); data[0] == 0x47 {
		return nil
	}

	in.Seek(96, os.SEEK_SET)
	in.Read(data)
	for i := 0; i < len(data); i++ {
		data[i] ^= 0xff
	}
	return data
}

// DefaultXORKey DefaultXORKey
var DefaultXORKey = []byte{
	0x59, 0x3d, 0x94, 0x69, 0x28, 0xed, 0xff, 0xe3, 0xc3, 0x6f, 0xac, 0xc5, 0xec, 0x52, 0x2b, 0xa3,
	0x79, 0x4e, 0xf5, 0x32, 0x43, 0x75, 0x88, 0x03, 0x1a, 0x84, 0x34, 0xcc, 0xb6, 0x53, 0x0d, 0x92,
	0x15, 0xd7, 0x2f, 0x30, 0xbd, 0x60, 0xb5, 0x17, 0x01, 0x9e, 0xdb, 0xb8, 0x56, 0x70, 0x33, 0x54,
	0x2e, 0x3b, 0xbf, 0x6a, 0x8c, 0x04, 0x41, 0xad, 0x6d, 0xf8, 0x58, 0x35, 0x98, 0x99, 0x24, 0x73,
	0x25, 0xf2, 0xb1, 0x5a, 0xb3, 0xc4, 0x8f, 0xd9, 0xef, 0xfb, 0x45, 0xc1, 0x37, 0x2a, 0x93, 0x4c,
	0x86, 0xda, 0x09, 0xae, 0x8d, 0xd8, 0x8a, 0x81, 0x7c, 0x44, 0x6b, 0xea, 0xf9, 0x66, 0x40, 0xa9,
	0xb9, 0x7a, 0x38, 0xe5, 0x29, 0xfc, 0x7f, 0x12, 0x4f, 0xcd, 0xba, 0xc7, 0x6c, 0xd6, 0xa0, 0x10,
	0xe2, 0x5e, 0xf0, 0xd1, 0xfd, 0xbb, 0xa6, 0x63, 0x05, 0xd5, 0x22, 0x9b, 0x9a, 0x00, 0xd3, 0x61,
	0x48, 0xaf, 0xee, 0xa7, 0x46, 0x77, 0x1f, 0x71, 0xde, 0x02, 0x42, 0x9c, 0xa2, 0xc2, 0xcf, 0xce,
	0x9d, 0x64, 0xf6, 0xca, 0xab, 0x14, 0x36, 0x0b, 0xd0, 0xdc, 0x3e, 0x7e, 0x0e, 0x72, 0x5f, 0x20,
	0x49, 0xa5, 0xfe, 0x23, 0xf4, 0x51, 0x95, 0x89, 0x87, 0xb0, 0x5d, 0x0f, 0x2c, 0x39, 0xa4, 0x5c,
	0xf1, 0x13, 0xcb, 0x57, 0x06, 0xf3, 0xeb, 0x97, 0x0c, 0x18, 0xb7, 0x21, 0xbc, 0x90, 0xe0, 0xb2,
	0x96, 0x1b, 0x27, 0xe9, 0x74, 0x19, 0x67, 0x6e, 0xc9, 0x55, 0xc0, 0x9f, 0xe8, 0xbe, 0xd4, 0x7b,
	0x83, 0x16, 0xd2, 0x2d, 0x4a, 0x0a, 0x7d, 0xb4, 0x82, 0x4d, 0xdd, 0x85, 0x3c, 0xe6, 0x50, 0x4b,
	0xc6, 0x80, 0xf7, 0x1d, 0xe7, 0x76, 0x3f, 0xfa, 0xe1, 0x78, 0xa8, 0x68, 0xa1, 0x1c, 0x91, 0xdf,
	0xaa, 0x3a, 0x26, 0x08, 0x8e, 0x62, 0x47, 0x5b, 0x1e, 0x65, 0xc8, 0x07, 0x11, 0x8b, 0x31, 0xe4,
}

// XORAllFileData XORAllFileData
func XORAllFileData(rootDir string, key []byte) {
	Debug("[XORAllFileData]%s", rootDir)
	fileInfos, _ := ioutil.ReadDir(rootDir)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			XORAllFileData(filepath.Join(rootDir, fileInfo.Name()), key)
		} else {
			filePath := filepath.Join(rootDir, fileInfo.Name())
			if file, err := os.OpenFile(filePath, os.O_RDWR, os.ModePerm); err == nil {
				keyLen := len(key)
				buff := make([]byte, 1024)
				keyIndex := 0
				for {
					n, err := file.Read(buff)
					if err != nil {
						break
					}
					data := buff[:n]
					for i, d := range data {
						data[i] = d ^ key[keyIndex%keyLen]
						keyIndex++
					}
					file.Seek(int64(-n), os.SEEK_CUR)
					file.Write(data)
				}
				file.Close()
			}
		}
	}
}

// GetMD5 GetMD5
func GetMD5(value string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(value)))
}

// AESDecryptCBC AESDecryptCBC
func AESDecryptCBC(encrypted []byte, key, iv []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	stream := cipher.NewCBCDecrypter(block, iv)
	stream.CryptBlocks(encrypted, encrypted)
	return encrypted[:len(encrypted)-int(encrypted[len(encrypted)-1])]
}

// CheckLiveURLState CheckLiveURLState
func CheckLiveURLState(url string) error {
	// connClient := core.NewConnClient()
	// if err := connClient.Start(url, "play"); err != nil {
	// 	// return err
	// 	return nil
	// }
	// connClient.Close(nil)
	return nil
}

// PermanentProtectCall PermanentProtectCall
func PermanentProtectCall(callback func()) {
	for {
		ProtectCall(callback)
	}
}

// ProtectCall ProtectCall
func ProtectCall(callback func()) {
	defer ProtectError()
	if callback != nil {
		callback()
	}
}

// DownloadFileToDisk DownloadFileToDisk
func DownloadFileToDisk(url string, savePath string, progressCallback func(current, total int64)) (size int64, err error) {
	if IsDebug() {
		Debug("[DownloadFileToDisk]%s << %s", savePath, url)
	} else {
		Info("[DownloadFileToDisk]%s", savePath)
	}

	var resp *http.Response
	mClinet := http.Client{Timeout: 20 * time.Second}
	if resp, err = mClinet.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	var file *os.File
	if file, err = os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, os.ModePerm); err != nil {
		return
	}
	defer file.Close()

	mTotalLength := resp.ContentLength
	mDownloadLength := int64(0)
	mProgressTime := time.Now()
	size = mTotalLength

	n := 0
	buff := make([]byte, 1024)
	for {
		if n, err = resp.Body.Read(buff); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		file.Write(buff[:n])
		mDownloadLength += int64(n)
		if progressCallback != nil {
			mNow := time.Now()
			if mNow.Sub(mProgressTime) >= time.Second {
				mProgressTime = mNow
				progressCallback(mDownloadLength, mTotalLength)
			}
		}
	}
	return
}

// LoadConfig LoadConfig
func LoadConfig(callback func() interface{}) error {
	if mConfig := callback(); mConfig != nil {
		//获取配置文件路径
		mConfigFile := os.Args[0]
		if index := strings.LastIndex(mConfigFile, "."); index > 0 {
			mConfigFile = mConfigFile[:index]
		}
		mConfigFile += ".cfg"
		//读取配置文件
		if data, err := ioutil.ReadFile(mConfigFile); err == nil {
			json.Unmarshal(data, mConfig)
		}
		//回写配置文件
		if data, err := json.MarshalIndent(mConfig, "", "\t"); err == nil {
			ioutil.WriteFile(mConfigFile, data, os.ModePerm)
		}
	}
	return nil
}
