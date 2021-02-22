package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

//查询视频源信息
func VideoAlter(startID, endID, mark int) {
	defer wg.Done() // goroutine结束就登记-1
	OldDB, err := InitDB(DsnOld)
	if err != nil {
		fmt.Println("OldDB InitDB failed! err:", err)
		return
	}
	NewDB, err := InitDB(DsnNew)
	if err != nil {
		fmt.Println("NewDB InitDB failed! err:", err)
		return
	}
	sqlOldInfos := "select id,user_id,title,video_url,video_times,charge,view_count,zan_count," +
		"share_count,status,remark_names,create_time,video_length from small_video where  ?<= id and id <?  "
	OldRows, err := OldDB.Query(sqlOldInfos, startID, endID)
	if err != nil {
		fmt.Println("db query failed! err:", err)
		return
	}
	//关闭rows 释放持有的数据库连接
	defer OldRows.Close()
	defer OldDB.Close()
	defer NewDB.Close()
	var OldVideoInfos NewVideoInfo
	for OldRows.Next() {
		OldVideoInfos.NewVideoCreator = "reptile"
		err = OldRows.Scan(
			&OldVideoInfos.NewVideoId,
			&OldVideoInfos.NewVideoUploadUser,
			&OldVideoInfos.NewVideoTitle,
			&OldVideoInfos.NewVideoUrl,
			&OldVideoInfos.NewVideoDuration,
			&OldVideoInfos.NewVideoPayCoin,
			&OldVideoInfos.NewVideoPlayCount,
			&OldVideoInfos.NewVideoPraiseCount,
			&OldVideoInfos.NewVideoShareCount,
			&OldVideoInfos.NewVideoStatus,
			&OldVideoInfos.NewVideoTags,
			&OldVideoInfos.NewVideoCreateTime,
			&OldVideoInfos.NewVideoByteSize)
		if err != nil {
			fmt.Println("rows scan failed!  err:", err)
			return
		}
		Now := time.Now().Format("2006-01-02-15:04:05")
		//拆标签字符串 转为json
		tagSlice := strings.Split(OldVideoInfos.NewVideoTags, ",")
		OldVideoInfos.TagsToJson(tagSlice, NewDB)
		fmt.Printf("[%v][Goroutine %v] video_id:%v \n", Now, mark, OldVideoInfos.NewVideoId)
		fmt.Printf("[%v][Goroutine %v] tags:%s \n", Now, mark, OldVideoInfos.NewVideoTags)
		OldVideoInfos.InsertVideoInfos(NewDB)
		//修改短视频资源文件夹名字
		OldVideoInfos.RenameDir(Now, mark)
	}

}

//视频信息插入新表
func (n *NewVideoInfo) InsertVideoInfos(NewDB *sql.DB) {
	sqlInsertNew := "insert into short_video_video_info (" +
		"video_id,video_updater,video_title," +
		"video_duration,video_pay_coin,video_play_count," +
		"video_praise_count,video_share_count,video_status," +
		"video_tags,video_creator,video_create_time,video_byte_size) " +
		"VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := NewDB.Exec(
		sqlInsertNew,
		n.NewVideoId,
		n.NewVideoUploadUser,
		n.NewVideoTitle,
		n.NewVideoDuration,
		n.NewVideoPayCoin,
		n.NewVideoPlayCount,
		n.NewVideoPraiseCount,
		n.NewVideoShareCount,
		n.NewVideoStatus,
		n.NewVideoTags,
		n.NewVideoCreator,
		n.NewVideoCreateTime,
		n.NewVideoByteSize)
	if err != nil {
		fmt.Println("NewDB Exec failed! err:", err)
		return
	}
}

//标签转json
func (n *NewVideoInfo) TagsToJson(tagSlice []string, NewDB *sql.DB) {
	var tagID int64
	var tag Tag
	n.NewVideoTagsJson = n.NewVideoTagsJson[0:0]
	for _, v := range tagSlice {
		//搜索新表中是否存在 当前标签
		sqlSelectTags := "select tag_id from short_video_tag_info where tag_name=?"
		_ = NewDB.QueryRow(sqlSelectTags, v).Scan(&tagID)
		tag.TagID = tagID
		tag.TagName = v
		tagsJson, err := json.Marshal(tag)
		if err != nil {
			fmt.Println("json Marshal failed! err:", err)
			return
		}
		InsertVideoRelations(tag.TagID, n.NewVideoId, NewDB)
		tagsJson = append(tagsJson, ',')
		n.NewVideoTagsJson = append(n.NewVideoTagsJson, tagsJson...)
		TagsStr := string(n.NewVideoTagsJson)
		TagsStr = strings.Trim(TagsStr, ",")
		TagsStr = "[" + TagsStr + "]"
		n.NewVideoTags = TagsStr
	}
}

//插入视频标签关系表
func InsertVideoRelations(tagId int64, videoId int, NewDB *sql.DB) {
	sqlInsertTags := "insert into short_video_tag_relation (tag_id,video_id) values(?,?)"
	_, err := NewDB.Exec(sqlInsertTags, tagId, videoId)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//修改资源文件夹名称
func (n *NewVideoInfo) RenameDir(Now string, mark int) {
	OldDirName := (strings.Split((strings.Split(n.NewVideoUrl, "smallvideos/"))[1], "/output.m3u8"))[0]
	OldDirPath := path.Join(ResDirPath, OldDirName)
	videoIdStr := strconv.Itoa(n.NewVideoId)
	NewDirPath := path.Join(ResDirPath, videoIdStr+"-A")
	fmt.Printf("[%v][Goroutine %v] OldDirPath:%v \n", Now, mark, OldDirPath)
	fmt.Printf("[%v][Goroutine %v] NewDirPath:%v \n", Now, mark, NewDirPath)
	err := os.Rename(OldDirPath, NewDirPath)
	if err != nil {
		fmt.Println("VideoDir Rename Failed! err:", err)
		return
	}
	XORAllFileData(NewDirPath, DefaultXORKey)
}
func main() {
	wg.Add(1) // 启动一个goroutine就登记+1
	go VideoAlter(63, 20000, 1)
	wg.Add(1) // 启动一个goroutine就登记+1
	go VideoAlter(20000, 40000, 2)
	wg.Add(1) // 启动一个goroutine就登记+1
	go VideoAlter(40000, 60000, 3)
	wg.Add(1) // 启动一个goroutine就登记+1
	go VideoAlter(60000, 80000, 4)
	wg.Add(1) // 启动一个goroutine就登记+1
	go VideoAlter(80000, 100000, 5)
	wg.Add(1) // 启动一个goroutine就登记+1
	go VideoAlter(100000, 130000, 6)
	wg.Wait()
}
