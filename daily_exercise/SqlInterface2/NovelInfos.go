package main

import "database/sql"

type NovelInfo struct {
	NovelOnlyQuery         bool               `json:"onlyQuery"`    // true 只查询 false 插入数据库
	NovelTitle             string             `json:"title"`        //小说标题
	NovelAuthor            string             `json:"author"`       //作者
	NovelIntroduce         string             `json:"introduce"`    //简介
	NovelWordTotalCount    int                `json:"totalWordNum"` //总字数
	NovelByteSize          int                `json:"totalSize"`    //总字节
	NovelChapterTotalCount int                `json:"totalNum"`     //当前总章节数
	NovelType              int                `json:"novelType"`    //小说类型|1-文本小说|2-有声小说
	NovelClassifyNames     string             `json:"type"`         //分类名称
	NovelCategoryName      string             `json:"categoryName"` //类别名称
	NovelNovelStatus       int                `json:"novelStatus"`  //状态（0完结，1连载）
	NovelChapter           []NovelChapterInfo `json:"chapters"`     //章节详情
	NovelID                int                //小说Id							ok
	NovelChapterFreeCount  int                //免费阅读章节数						ok
	NovelCoinPrice         int                //收费金币数						ok
	NovelPayType           int                //收费类型（0免费，1VIP/金币 2付费）	ok
	NovelNovelEnable       int                //是否启用（0禁用，1启用）			ok
	NovelNovelPopularity   int                //人气值							ok
	NovelNovelScore        float32            //评分								ok
	NovelBaseReadCount     int                //在读人数初始值						ok
	NovelSearchCount       int                //搜索数							ok
	NovelFavoriteCount     int                //初始收藏数
	NovelWordCountType     int                //当前总字数分类 1.100万字以下 2.100万~200万 3.200万字以上	ok
	NovelCreateTime        string             //创建时间							ok
	NovelCreator           string             //小说创建者						ok
	NovelCategoryID        int64              //类别Id							ok
	NovelChapterToSlice    []int              //小说所有章节的Id 切片
	NovelChapterToJson     string             //小说所有章节ID 的json string
	NovelClassifyToJson    []NovelClassify    //小说分类
	NovelClassifyToString  string             //小说分类 json
	NovelUpdateTime        string             //小说更新时间
	NovelMachineID         int                //小说资源服务器iD
	NvlDBTx                *sql.Tx
}
type NovelChapterInfo struct {
	NovelChapterByteSize int    `json:"chapterSize"` //章节大小 byte
	NovelChapterName     string `json:"subTitle"`    //章节名称
	NovelChapterSort     int    `json:"chapterID"`   //章节排序
	NovelChapterTime     int    `json:"chapterTime"` //有声小说时长
}
type NovelClassify struct {
	NovelCategoryId   int64  `json:"categoryId"`   //类别ID
	NovelClassifyID   int64  `json:"classifyId"`   //分类Id
	NovelClassifyName string `json:"classifyName"` //分类名称
}
