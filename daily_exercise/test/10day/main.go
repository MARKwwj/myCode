package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type videoInfo struct {
	videoID    int
	videoTimes int
	videoUrl   string
}
type allVideoInfos struct {
	Info map[int]*videoInfo
}

// 定义一个初始化数据库的函数
func initDB() (err error) {
	dsn := "root:ZsNice2020.@tcp(199.180.114.169:6033)/jtest"
	//不会校验账号密码是否正确
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	//尝试与数据库建立连接（检验dsn是否正确）
	err = db.Ping()
	if err != nil {
		fmt.Println("与数据库连接失败！")
		return err
	}
	return nil
}
func QueryRow() {
	sqlStr := "select id,video_times,video_url from small_video where id = ?"
	var v videoInfo
	//非常重要 确保queryRow 之后使用Scan方法 否则持有的数据库连接将不会被释放
	err := db.QueryRow(sqlStr, 1).Scan(&v.videoID, &v.videoTimes, &v.videoUrl)
	if err != nil {
		fmt.Println("scan failed! err:", err)
		return
	}
	fmt.Println(v.videoID, " ", v.videoTimes, " ", v.videoUrl)
}
func QueryMultiRow() {
	sqlStr := "select id,video_times,video_url from small_video where video_times<?"
	//非常重要 确保queryRow 之后使用Scan方法 否则持有的数据库连接将不会被释放
	rows, err := db.Query(sqlStr, 100)
	if err != nil {
		fmt.Println("scan failed! err:", err)
		return
	}
	//非常重要 关闭rows 释放持有的数据库连接
	defer rows.Close()
	a := &allVideoInfos{make(map[int]*videoInfo, 3000)}
	var curInfo *videoInfo
	for rows.Next() {
		var videoID int
		var videoTimes int
		var videoUrl string
		err := rows.Scan(&videoID, &videoTimes, &videoUrl)
		if err != nil {
			fmt.Println("rows Scan Failed! err:", err)
			return
		}
		curInfo = &videoInfo{
			videoID,
			videoTimes,
			videoUrl,
		}
		a.Info[videoID] = curInfo
		fmt.Println(curInfo)
	}
}
func main() {
	err := initDB()
	if err != nil {
		fmt.Println("initDB failed! err:", err)
		return
	}
	QueryMultiRow()
}
