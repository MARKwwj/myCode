package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var coverName []string

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

func checkCoverSize(NovelPath string) {
	fileObj, err := ioutil.ReadDir(NovelPath)
	if err != nil {
		fmt.Println("ioutil.ReadDir(NovelPath) failed,err:", err)
		return
	}
	var k int
	for _, file := range fileObj {
		if strings.Contains(file.Name(), ".zip") {
			continue
		}
		curCoverPath := path.Join(NovelPath, file.Name(), "cover.webp")
		curJpgCoverPath := path.Join(NovelPath, file.Name(), "cover.jpg")

		var coverObj os.FileInfo
		coverObj, err = os.Stat(curCoverPath)
		if err != nil {
			fmt.Println(" os.Stat(curCoverPath) failed,err:", err)
			continue
		}
		_ = os.Remove(curJpgCoverPath)

		if coverObj.Size() == 0 {
			fmt.Println(curCoverPath)
			_ = os.Remove(curCoverPath)
			srcNewCover := coverName[k]
			fmt.Println("************************************************")
			fmt.Printf("coverUsedNums:%v \n", k+1)
			fmt.Printf("srcNewCover:%v \n", srcNewCover)
			fmt.Printf("curCoverPath:%v \n", curCoverPath)
			fmt.Println("************************************************")
			CopyFile(srcNewCover, curCoverPath)
		}

	}
}

func main() {
	NovelPath := "/home/datadrive/resources/res3/novel/crypted/chapterDetail/"
	srcFilePath := "/home/datadrive/resources/res3/novel/crypted/coverWebp/"
	//NovelPath := "D:\\desktop\\test"
	//srcFilePath := "D:\\desktop\\ppp"

	getNewCover(srcFilePath)
	checkCoverSize(NovelPath)

}
