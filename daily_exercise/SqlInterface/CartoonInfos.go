package main

type CartoonInfo struct {
	CartoonOnlyQuery         bool                 `json:"onlyQuery"`    // true 只查询 false 插入数据库
	CartoonTitle             string               `json:"title"`        //漫画标题
	CartoonIntroduce         string               `json:"introduce"`    //简介
	CartoonChapterTotalCount int                  `json:"totalNum"`     //当前总章节数
	CartoonClassifyNames     string               `json:"type"`         //漫画分类
	CartoonCategoryName      string               `json:"categoryName"` //漫画类别
	CartoonChapterInfos      []CartoonChapterInfo `json:"chapters"`     //章节详情
	CartoonID                int64                //漫画Id				ok
	CartoonAuthor            string               //作者 来自于网络		ok
	CartoonChapterFreeCount  int                  //免费阅读章节数			ok
	CartoonCoinPrice         int                  //收费金币数			ok
	CartoonStatus            int                  //状态（0完结，1连载）	ok
	CartoonPayType           int                  //收费类型（0免费，1付费）ok
	CartoonEnable            int                  //是否启用（0禁用，1启用）ok
	CartoonPopularity        int                  //人气值				ok
	CartoonScore             float32              //评分					ok
	CartoonFavoriteCount     int                  //收藏数				ok
	CartoonCreateTime        string               //创建时间				ok
	CartoonCreator           string               //漫画创建者 reptile	ok
	CartoonCategoryID        int                  //漫画类别ID			ok
	CartoonChapterToSlice    []int                //漫画所有章节的Id 切片
	CartoonChapterToJson     string               //漫画所有章节的Id json数组
	CartoonClassifyToJson    []CartoonClassify    //漫画分类
	CartoonClassifyToString  string               //漫画分类 json
}
type CartoonChapterInfo struct {
	CartoonChapterSort       int              `json:"chapterID"`   //排序 章节Id
	CartoonChapterName       string           `json:"subTitle"`    //章节标题
	CartoonChapterByteSize   int              `json:"chapterSize"` //漫画字节大小
	CartoonPicInfos          []CartoonPicInfo `json:"chapterPics"` //漫画图片详情
	CartoonPicCount          int              //当前章节漫画图片数量
	CartoonChapterPicsSlice  []int            //章节所包含的pic_id
	CartoonChapterPicsString string           //章节所包含的pic_id json数组
}

type CartoonPicInfo struct {
	CartoonChapterID int `json:"chapterID"` //章节Id
	CartoonPicNum    int `json:"picIndex"`  //图片标识数
	CartoonHeight    int `json:"height"`    //图片高
	CartoonWidth     int `json:"width"`     //图片宽
}

type CartoonClassify struct {
	CartoonClassifyID   int64  `json:"classify_id"`   //分类Id
	CartoonClassifyName string `json:"classify_name"` //分类名称
}
