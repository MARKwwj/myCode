package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type txtNovelInfo struct {
	title   string
	novelID int
}
type infos struct {
	txtInfo map[int]*txtNovelInfo
}

var coverName []string
var TxtIDSlice []int

func getNovelID() *infos {
	DB, err := sql.Open("mysql", "root:yxPvqJlbYBRrIs0z@tcp(110.92.66.88:6033)/res_novel_db")
	if err != nil {
		fmt.Println("sql open failed! err:", err)
		return nil
	}
	defer DB.Close()
	err = DB.Ping()
	if err != nil {
		fmt.Println("DB ping failed! err:", err)
		return nil
	}
	rows, err := DB.Query("select title,novel_id from novel_info where novel_type=?", 1)
	if err != nil {
		fmt.Println("DB query failed! err:", err)
		return nil
	}
	var novelID int
	var title string
	txtInfos := make(map[int]*txtNovelInfo, 1000)
	for rows.Next() {
		err = rows.Scan(&title, &novelID)
		if err != nil {
			fmt.Println("rows Scan faild,err", err)
			return nil
		}
		newTxtInfo := &txtNovelInfo{
			title:   title,
			novelID: novelID,
		}
		txtInfos[novelID] = newTxtInfo
	}
	newInfos := infos{
		txtInfo: txtInfos,
	}
	return &newInfos
}
func CopyFile(srcFileName string, dstFileName string) {
	//打开源文件
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		log.Fatalf("源文件读取失败,原因是:%v\n", err)
	}
	defer func() {
		err = srcFile.Close()
		if err != nil {
			log.Fatalf("源文件关闭失败,原因是:%v\n", err)
		}
	}()

	//创建目标文件,稍后会向这个目标文件写入拷贝内容
	distFile, err := os.Create(dstFileName)
	if err != nil {
		log.Fatalf("目标文件创建失败,原因是:%v\n", err)
	}
	defer func() {
		err = distFile.Close()
		if err != nil {
			log.Fatalf("目标文件关闭失败,原因是:%v\n", err)
		}
	}()
	//定义指定长度的字节切片,每次最多读取指定长度
	var tmp = make([]byte, 1024*4)
	//循环读取并写入
	for {
		n, err := srcFile.Read(tmp)
		n, _ = distFile.Write(tmp[:n])
		if err != nil {
			if err == io.EOF { //读到了文件末尾,并且写入完毕,任务完成返回(关闭文件的操作由defer来完成)
				return
			} else {
				log.Fatalf("拷贝过程中发生错误,错误原因为:%v\n", err)
			}
		}
	}
}

func decode(filePath string) {
	var key = []byte{
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
func alterNanNvCover() {
	newFilePath := "/home/datadrive/resources/res3/novel/crypted/cover_new"
	NovelPath := "/home/datadrive/resources/res3/novel/crypted/chapterDetail/"
	info := getNovelID()
	fielObj, err := ioutil.ReadDir(newFilePath)
	if err != nil {
		fmt.Println("ioutil ReadDir Failed,err", err)
		return
	}
	for _, v := range info.txtInfo {
		for _, file := range fielObj {
			fileSplitJpgName := strings.Trim(file.Name(), ".jpg")
			if v.title == fileSplitJpgName {
				curJpgCoverPath := path.Join(NovelPath, strconv.Itoa(v.novelID), "cover.jpg")
				newJpgCoverPath := path.Join(newFilePath, file.Name())
				fmt.Println(curJpgCoverPath)
				fmt.Println(newJpgCoverPath)
				err = os.Remove(curJpgCoverPath)
				if err != nil {
					fmt.Println("os remove failed,err", err)
				}
				CopyFile(newJpgCoverPath, curJpgCoverPath)
			}
		}
	}
}
func ConvertToWebp(src, dst string) {
	//var out bytes.Buffer
	ffmpegPath, _ := exec.LookPath("ffmpeg")
	ffmpegPath, _ = filepath.Abs(ffmpegPath)
	cmd := exec.Command(ffmpegPath, "-i", src, dst)
	err := cmd.Start()
	if err != nil {
		fmt.Println("cmd.Start failed!,err:", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("cmd.Wait() failed! err:", err)
		return
	}
}

func main() {
	NovelPath := "D:\\desktop\\aaa\\test"
	backupsCover := "D:\\desktop\\aaa\\covers"
	//NovelPath := "/home/datadrive/resources/res3/novel/crypted/chapterDetail/"
	//backupsCover := "/home/datadrive/resources/res3/novel/crypted/covers"

	fileObj, err := ioutil.ReadDir(NovelPath)
	if err != nil {
		fmt.Println("ioutil ReadDir failed,err", err)
		return
	}

	for _, file := range fileObj {
		if !file.IsDir() {
			continue
		}

		curCoverPath := path.Join(NovelPath, file.Name(), "cover.jpg")
		newCoverPath := path.Join(NovelPath, file.Name(), "cover.webp")

		_, err = os.Stat(newCoverPath)
		if err != nil {

			_, err = os.Stat(curCoverPath)
			if err != nil {

				var coverObj []os.FileInfo
				coverObj, err = ioutil.ReadDir(backupsCover)
				if err != nil {
					fmt.Println("ioutil ReadDir failed,err", err)
					return
				}
				for _, v := range coverObj {
					if strings.Contains(v.Name(), "-used") {
						continue
					} else {
						backupCoverPath := path.Join(backupsCover, v.Name())
						CopyFile(backupCoverPath, newCoverPath)
						err = os.Rename(backupCoverPath, backupCoverPath+"-used")
						if err != nil {
							fmt.Println("err-backupCoverPath:", backupCoverPath)
							fmt.Println("os Rename backcovername failed,err:", err)
						}
						break
					}
				}
				continue
			}
			fmt.Println("*************************************************")
			fmt.Println(curCoverPath)
			fmt.Println(newCoverPath)
			decode(curCoverPath)
			ConvertToWebp(curCoverPath, newCoverPath)
			//_ = os.Remove(curCoverPath)
			decode(newCoverPath)
			decode(curCoverPath)
			fmt.Println("*************************************************")
		}
	}
}
