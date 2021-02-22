package main

import (
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

//video new information
type NewVideoInfo struct {
	NewVideoId          int    //视频ID
	NewVideoUploadUser  int    //视频上传用户id
	NewVideoDuration    int    //视频时长
	NewVideoPayCoin     int    //支付金币
	NewVideoPlayCount   int    //播放数
	NewVideoPraiseCount int    //点赞数
	NewVideoShareCount  int    //分享数
	NewVideoStatus      int    //视频状态
	NewVideoByteSize    int    //字节大小
	NewVideoTitle       string //视频标题
	NewVideoTags        string //视频标签 json
	NewVideoCreator     string //创建者
	NewVideoCreateTime  string //创建时间
	NewVideoUrl         string //视频m3u8文件路径
	NewVideoTagsJson    []byte //视频标签 json
}
type AllInfos struct {
	infos map[int]*NewVideoInfo
}
type Tag struct {
	TagID   int64  `json:"tag_id"`
	TagName string `json:"tag_name"`
}

//dataSourceName
const DsnOld = "root:ZsNice2020.@tcp(199.180.114.169:6033)/jtest"
const DsnNew = "root:ZsNice2020.@tcp(58.82.232.37:6033)/res_short_video_db"

//const DsnNew = "root:ZsNice2020.@tcp(199.180.114.169:6033)/res_short_video_db"

const ResDirPath = "/home/datadrive/resources/res2/smallvideo/crypted/videos"

//const ResDirPath = "D:\\desktop\\r1"

var wg sync.WaitGroup
