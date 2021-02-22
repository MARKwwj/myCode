package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path"
	"strconv"
	"time"
)

var DB *sql.DB

//const NovelPath = "/home/datadrive/resources/res3/novel/crypted/CapturerNovel_v21.0121/Data"
//const SoundNovelPath = "/home/datadrive/resources/res3/novel/crypted/CapturerSoundNovel_v21.0122/Data"
//const NovelPath = "D:\\desktop"
//const SoundNovelPath = "D:\\desktop"
const NovelDsnPro = "root:yxPvqJlbYBRrIs0z@tcp(110.92.66.88:6033)/res_novel_db"
const NovelPath = "/home/datadrive/resources/res3/novel/crypted/chapterDetail"


func initDB() (err error) {
	DB, err = sql.Open("mysql", NovelDsnPro)
	if err != nil {
		return fmt.Errorf("sql open failed %v", err)
	}
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("DB ping failed! %v", err)
	}
	DB.SetConnMaxLifetime(time.Second * 500)
	DB.SetMaxOpenConns(10)
	DB.SetConnMaxIdleTime(5)
	return nil
}
func alter(novelType int, filePath string) {
	err := initDB()
	if err != nil {
		fmt.Println("initDB failed ! err:", err)
		return
	}
	rows, err := DB.Query("select novel_id from novel_info where novel_type =?", novelType)
	if err != nil {
		fmt.Println("row Scan failed! err:", err)
		return
	}
	var ID int
	var IDSlice []int
	for rows.Next() {
		_ = rows.Scan(&ID)
		IDSlice = append(IDSlice, ID)
	}
	rows.Close()
	for _, v := range IDSlice {
		curFilePath := path.Join(filePath, strconv.Itoa(v))
		_, err = os.Stat(curFilePath)
		if err != nil {
			fmt.Println(v)
		}
	}
}
func main() {
	alter(1, NovelPath)
	alter(2, NovelPath)
	//a := []byte{
	//	0xe7, 0x8f, 0xa0, 0xe8, 0x83, 0x8e, 0xe6, 0x9a, 0x97, 0xe7, 0xbb, 0x93, 0x3a, 0xe7, 0x8e, 0x8b,
	//	0xe6, 0xb5, 0xa9,
	//}
	//fmt.Printf("%s", string(a))
}
