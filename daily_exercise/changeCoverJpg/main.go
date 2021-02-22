package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

var SoundIDSlice []int
var TxtIDSlice []int
var coverName []string

func getNovelID() {
	DB, err := sql.Open("mysql", "root:yxPvqJlbYBRrIs0z@tcp(110.92.66.88:6033)/res_novel_db")
	if err != nil {
		fmt.Println("sql open failed! err:", err)
		return
	}
	err = DB.Ping()
	if err != nil {
		fmt.Println("DB ping failed! err:", err)
		return
	}
	//rows, err := DB.Query("select novel_id from novel_info where novel_type=?", 2)
	//if err != nil {
	//	fmt.Println("DB query failed! err:", err)
	//	return
	//}
	//var SoundNovelID int
	//for rows.Next() {
	//	err = rows.Scan(&SoundNovelID)
	//	if err != nil {
	//		fmt.Println("rows Scan failed! err:", err)
	//		return
	//	}
	//	SoundIDSlice = append(SoundIDSlice, SoundNovelID)
	//}
	txtRows, err := DB.Query("select novel_id from novel_info where novel_type=?", 2)
	if err != nil {
		fmt.Println("DB query failed! err:", err)
		return
	}
	var TxtNovelID int
	for txtRows.Next() {
		err = txtRows.Scan(&TxtNovelID)
		if err != nil {
			fmt.Println("rows Scan failed! err:", err)
			return
		}
		TxtIDSlice = append(TxtIDSlice, TxtNovelID)
	}
}
func coverChange(filePath string) {
	for k, v := range TxtIDSlice {
		NovelPath := path.Join(filePath, strconv.Itoa(v))
		fmt.Println(NovelPath)
		curJpgCoverPath := path.Join(NovelPath, "cover.jpg")
		cruWebpCoverPath := path.Join(NovelPath, "cover.webp")
		srcCoverPath := coverName[k]
		_ = os.Remove(curJpgCoverPath)
		_ = os.Remove(cruWebpCoverPath)
		CopyFile(srcCoverPath, cruWebpCoverPath)
	}
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
func getNewCover(srcFilePath string) {
	srcFileObj, err := ioutil.ReadDir(srcFilePath)
	if err != nil {
		fmt.Println("ioutil readDir failed! err:", err)
		return
	}
	for _, src := range srcFileObj {
		curCoverName := path.Join(srcFilePath, src.Name())
		coverName = append(coverName, curCoverName)
	}
}

//func txtChangeCover(txtNovelPath string) {
//	for
//}

func main() {
	filePath := "/home/datadrive/resources/res3/novel/crypted/chapterDetail/"
	srcFilePath := "/home/datadrive/resources/res3/novel/crypted/webpRecord/"
	//filePath := "D:\\desktop\\cover_rep\\222"
	//srcFilePath := "D:\\desktop\\cover_rep\\111"
	TxtIDSlice = []int{2750,2732,2731,2724,2641,2609,2782,2908,2872,2815,2945,2935,2922}
	//getNovelID()
	getNewCover(srcFilePath)
	coverChange(filePath)
}
