package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type NovelChapterInfos struct {
	NovelID              int
	NovelChaptersIDStr   string
	NovelChaptersIDSlice []int
	NovelIDSlice         []int
}

//type NovelChapters struct {
//	NovelChapterID       int
//	NovelID              int
//	NovelChapterName     string
//	NovelChapterNo       int
//	NovelChapterByteSize int
//	NovelChapterTime     int
//}
func InitDB(Dsn string) (DB *sql.DB) {
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
func FFmpegSoundNovelToTs(NovelPath string) {
	//SoundNovel Slice
	fileInfos, err := ioutil.ReadDir(NovelPath)
	if err != nil {
		fmt.Println("file Read Failed! err:", err)
		return
	}
	//先解密
	//XORAllFileData(NovelPath, DefaultXORKey)
	for _, v := range fileInfos {
		//不是文件夹 且不是jpg结尾
		if v.IsDir() == false && strings.HasSuffix(v.Name(), "jpg") == false {
			Name := (strings.Split(v.Name(), "."))[0]
			DirName := path.Join(NovelPath, v.Name())
			DirPath := path.Join(NovelPath, Name)
			_, err = os.Stat(DirPath)
			switch {
			case os.IsNotExist(err):
				_ = os.Mkdir(DirPath, 777)
			case os.IsExist(err):
				_ = os.Mkdir(DirPath+"-A", 777)
			case err != nil:
				fmt.Println("os.Stat DirPath  Failed! err:", err)
				return
			}
			//再切片
			ConvertToTs(DirPath, DirName)
		}
	}
}
func ConvertToTs(ToFilePath, mp3FilePath string) {
	//var out bytes.Buffer
	ffmpegPath, _ := exec.LookPath("ffmpeg")
	ffmpegPath, _ = filepath.Abs(ffmpegPath)
	cmd := exec.Command(ffmpegPath, "-i", mp3FilePath, "-acodec", "copy", "-f", "segment", "-segment_time", "10", "-segment_list",
		ToFilePath+"/output.m3u8", "-segment_format", "mpegts", ToFilePath+"/output%04d.ts")
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
func ChangeChapterRelationIDJson() {
	DB := InitDB("root:123456@tcp(192.168.100.51:3306)/novel_test")
	rows, err := DB.Query("select novel_id from novel_info")
	if err != nil {
		fmt.Println("DB Query Failed! err:", err)
		return
	}
	var Nvl NovelChapterInfos
	for rows.Next() {
		var NovelID int
		err = rows.Scan(&NovelID)
		if err != nil {
			fmt.Println("rows Scan Failed! err:", err)
			return
		}
		Nvl.NovelIDSlice = append(Nvl.NovelIDSlice, NovelID)
	}
	for _, v := range Nvl.NovelIDSlice {
		Nvl.NovelChaptersIDSlice = Nvl.NovelChaptersIDSlice[0:0]
		ChapterIDRows, err := DB.Query("select chapter_id from novel_chapter_info where novel_id = ?", v)
		if err != nil {
			fmt.Println("chapterID DB query Failed! err:", err)
			return
		}
		for ChapterIDRows.Next() {
			var ChapterID int
			err := ChapterIDRows.Scan(&ChapterID)
			if err != nil {
				fmt.Println("chaperIDRows Scan Failed! err:", err)
				return
			}
			Nvl.NovelChaptersIDSlice = append(Nvl.NovelChaptersIDSlice, ChapterID)
		}
		ChapterIDByte, _ := json.Marshal(Nvl.NovelChaptersIDSlice)
		Nvl.NovelChaptersIDStr = string(ChapterIDByte)
		_, err = DB.Exec("update novel_chapter_relation set novel_chapters=? where novel_id =?", Nvl.NovelChaptersIDStr, v)
		if err != nil {
			fmt.Println("DB Exec Failed! err:", err)
			return
		}
	}
}
//
//func main() {
//	NovelPath := "D:\\desktop\\bbbb"
//	fileObj, err := ioutil.ReadDir(NovelPath)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	for _, v := range fileObj {
//		curPath := path.Join(NovelPath, v.Name())
//		fmt.Println(curPath)
//		FFmpegSoundNovelToTs(curPath)
//	}
//	//ChangeChapterRelationIDJson()
//}
