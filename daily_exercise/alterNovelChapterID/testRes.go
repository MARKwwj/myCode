package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

func InitDB2(Dsn string) (DB *sql.DB) {
	//不会校验账号密码是否正确
	DB, err := sql.Open("mysql", Dsn)
	if err != nil {
		fmt.Println("sql open failed! err:", err)
		return nil
	}
	//尝试与数据库连接，并校验dsn是否正确
	err = DB.Ping()
	if err != nil {
		fmt.Println("DB ping failed! err:", err)
		return nil
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	return DB
}

func rename() {
	filePath := "D:\\desktop\\bbbb"
	fileObj, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Println("ioutil failed ")
		return
	}
	i := 100
	for _, v := range fileObj {
		i++
		curFilePath := path.Join(filePath, v.Name())
		newFilePath := path.Join(filePath, strconv.Itoa(i))
		err = os.Rename(curFilePath, newFilePath)
		if err != nil {
			fmt.Println("rename failed ! ")
		}
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
func RenameSound() {
	filePath := "D:\\desktop\\bbbb"
	fileObj, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Println("ioutil failed ")
		return
	}
	i := 100
	for _, v := range fileObj {
		i++
		OldJpg := path.Join(filePath, v.Name(), "cover.jpg")
		_ = os.Remove(OldJpg)
		NewJpgFilePath := path.Join("D:\\desktop\\Data_All\\tu\\1\\02", strconv.Itoa(i)+".jpg")
		CopyFile(NewJpgFilePath, OldJpg)
	}
}
func CopyJpg() {
	filePath := "D:\\desktop\\xxx"
	fileObj, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Println("ioutil failed ")
		return
	}
	i := 0
	for _, v := range fileObj {
		i++
		curFilePath := path.Join(filePath, v.Name())
		_ = os.Remove(path.Join(curFilePath, "cover.jpg"))
		jpgPath := path.Join("D:\\desktop\\Data_All\\tu\\1\\02", strconv.Itoa(i)+".jpg")
		CopyFile(jpgPath, path.Join(curFilePath, "cover.jpg"))
		//newFilePath := path.Join(filePath, strconv.Itoa(i))
		//_ = os.Rename(curFilePath, newFilePath)
	}
}
func renameTxt() {
	DB := InitDB2("root:ZsNice2020.@tcp(199.180.114.169:6033)/res_novel_db")
	filePath := "D:\\desktop\\bbbb"
	fileObj, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Println("ioutil failed ")
		return
	}
	for _, v := range fileObj {
		curDirPath := path.Join(filePath, v.Name())
		var byteM []byte
		id, _ := strconv.Atoi(v.Name())
		fmt.Println(id)
		err = DB.QueryRow("select novel_chapters from novel_chapter_relation where novel_id=?", id).Scan(&byteM)
		if err != nil {
			fmt.Println("DB QueryRow failed !err:", err)
			return
		}
		var slice []int
		_ = json.Unmarshal(byteM, &slice)
		i := 0
		for j := 0; j < 20; j++ {
			curTxtPath := path.Join(curDirPath, strconv.Itoa(j+1))
			NewTxtPath := path.Join(curDirPath, strconv.Itoa(slice[i]))
			err = os.Rename(curTxtPath, NewTxtPath)
			if err != nil {
				fmt.Println(err)
			}
			if i < 19 {
				i++
			}
		}
	}
}
func main() {
	//rename()
	//CopyJpg()
	renameTxt()
}
