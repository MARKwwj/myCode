package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func NovelMainFunc(c *gin.Context) {
	var NvlInfo NovelInfo
	if err := c.BindJSON(&NvlInfo); err != nil {
		fmt.Println(err)
	}
	// 先验证是否重复
	//1、文字小说  2、有声小说
	var ok bool
	var err error
	if NvlInfo.NovelType == 1 {
		ok, err = PreventRepeating(NvlInfo.NovelTitle, NovelTitleMd5)
		if err != nil {
			fmt.Println("PreventRepeating Failed! err:", err)
			return
		}
	} else if NvlInfo.NovelType == 2 {
		ok, err = PreventRepeating(NvlInfo.NovelTitle, NovelSoundTitleMd5)
		if err != nil {
			fmt.Println("PreventRepeating Failed! err:", err)
			return
		}
	}
	if NvlInfo.NovelOnlyQuery == true {
		// true 重复 false 未重复
		if ok {
			c.JSON(http.StatusOK, gin.H{"repeat": true})
			fmt.Println("repeat: true")
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"repeat": false})
			fmt.Println("repeat: false")
			return
		}
	} else {
		if ok {
			c.JSON(http.StatusOK, gin.H{"repeat": true})
			fmt.Println("repeat: true")
			return
		}
		NvlInfo.InsertIntoClassifyInfo()
		NvlInfo.SelectNovelCategoryID()
		NvlInfo.InsertIntoNovelInfo()
		NvlInfo.InsertIntoNovelChapterInfos()
		NvlInfo.InsertIntoNovelChapterRelation()
		NvlInfo.InsertIntoNovelClassifyRelation()
		NvlInfo.InsertIntoNovelCategoryRelation()
		if NvlInfo.NovelType == 1 {
			NewTitleWriteToMap(NvlInfo.NovelTitle, NovelTitleMd5)
		} else if NvlInfo.NovelType == 2 {
			NewTitleWriteToMap(NvlInfo.NovelTitle, NovelSoundTitleMd5)
		}
		c.JSON(http.StatusOK, gin.H{"novelID": NvlInfo.NovelID,"chapterID":NvlInfo.NovelChapterToSlice})
	}
}
func (n *NovelInfo) InsertIntoClassifyInfo() {
	var NvlClassify NovelClassify
	var TagsStr string
	if n.NovelClassifyNames == "" {
		n.NovelClassifyToString = "[]"
		return
	}
	tagSlice := strings.Split(n.NovelClassifyNames, ",")
	for _, ClassifyName := range tagSlice {
		var ClassifyID int64
		ClassifySql := "select classify_id from novel_classify_info where classify_name=?"
		err := NvlDB.QueryRow(ClassifySql, ClassifyName).Scan(&ClassifyID)
		switch {
		case err == sql.ErrNoRows:
			ClassifyInsertSql := "insert into novel_classify_info(category_id,classify_name) values(1,?)"
			result, err := NvlDB.Exec(ClassifyInsertSql, ClassifyName)
			if err != nil {
				fmt.Println("NvlDB Exec failed! err:", err)
				return
			}
			ClassifyID, _ = result.LastInsertId()
		case err != nil:
			fmt.Println("NvlDB queryRow failed! err:", err)
			return
		}
		//转 json
		NvlClassify.NovelClassifyID = ClassifyID
		NvlClassify.NovelClassifyName = ClassifyName
		n.NovelClassifyToJson = append(n.NovelClassifyToJson, NvlClassify)
	}
	tagsJson, _ := json.Marshal(n.NovelClassifyToJson)
	TagsStr = string(tagsJson)
	n.NovelClassifyToString = TagsStr
}
func (n *NovelInfo) InsertIntoNovelInfo() {
	if n.NovelAuthor == "" {
		n.NovelAuthor = "来自于网络"
	}
	if n.NovelIntroduce == "" {
		n.NovelIntroduce = n.NovelTitle
	}
	n.NovelPayType = 1
	n.NovelCoinPrice = 2
	n.NovelNovelEnable = 0
	n.NovelChapterFreeCount = 1
	n.NovelNovelPopularity = rand.Intn(87000) + 12900
	n.NovelNovelScore = (float32(rand.Intn(9)) / 10) + float32(8.8)
	n.NovelBaseReadCount = rand.Intn(8700) + 1290
	n.NovelSearchCount = rand.Intn(970) + 129
	n.NovelCreateTime = time.Now().Format("2006-01-02 15:04:05")
	n.NovelCreator = "reptile"
	////小说所有章节的Id json数组
	//for k, _ := range n.NovelChapter {
	//	n.NovelChapterToSlice = append(n.NovelChapterToSlice, k+1)
	//}
	//byteSlice, _ := json.Marshal(n.NovelChapterToSlice)
	//n.NovelChapterToJson = string(byteSlice)
	switch {
	case n.NovelWordTotalCount <= 1000000:
		n.NovelWordCountType = 1
	case n.NovelWordTotalCount < 2000000:
		n.NovelWordCountType = 2
	case n.NovelWordTotalCount >= 2000000:
		n.NovelWordCountType = 3
	}
	NovelInfoInsertSql := "INSERT INTO novel_info" +
		"(title, " +
		"author, " +
		"introduce, " +
		"chapter_total_count, " +
		"chapter_free_count, " +
		"coin_price, " +
		"status, " +
		"pay_type, " +
		"enable, " +
		"popularity, " +
		"score, " +
		"base_read_count, " +
		"novel_classifies, " +
		"search_count, " +
		"word_total_count, " +
		"word_count_type, " +
		"create_time, " +
		"novel_creator," +
		"novel_type," +
		"novel_byte_size) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := NvlDB.Exec(NovelInfoInsertSql,
		n.NovelTitle,
		n.NovelAuthor,
		n.NovelIntroduce,
		n.NovelChapterTotalCount,
		n.NovelChapterFreeCount,
		n.NovelCoinPrice,
		n.NovelNovelStatus,
		n.NovelPayType,
		n.NovelNovelEnable,
		n.NovelNovelPopularity,
		n.NovelNovelScore,
		n.NovelBaseReadCount,
		n.NovelClassifyToString,
		n.NovelSearchCount,
		n.NovelWordTotalCount,
		n.NovelWordCountType,
		n.NovelCreateTime,
		n.NovelCreator,
		n.NovelType,
		n.NovelByteSize)
	if err != nil {
		fmt.Println("NvlDB.Exec NovelInfoInsertSql Failed! err:", err)
		return
	}
	n.NovelID, err = result.LastInsertId()
	if err != nil {
		fmt.Println("NovelID  result.LastInsertId() Failed! err:", err)
		return
	}
}

func (n *NovelInfo) InsertIntoNovelChapterInfos() {
	for _, v := range n.NovelChapter {
		NovelChapterInfoInsertSql := "INSERT INTO novel_chapter_info(novel_id, chapter_name, chapter_no, chapter_byte_size,chapter_time) VALUES (?,?,?,?,?)"
		result, err := NvlDB.Exec(NovelChapterInfoInsertSql, n.NovelID, v.NovelChapterName, v.NovelChapterSort, v.NovelChapterByteSize, v.NovelChapterTime)
		if err != nil {
			fmt.Println("NvlDB.Exec NovelChapterInfoInsertSql Failed ! err:", err)
			return
		}
		chapterID, _ := result.LastInsertId()
		n.NovelChapterToSlice = append(n.NovelChapterToSlice, int(chapterID))
	}
	byteSlice, _ := json.Marshal(n.NovelChapterToSlice)
	n.NovelChapterToJson = string(byteSlice)
}

func (n *NovelInfo) InsertIntoNovelChapterRelation() {
	NovelChapterRelationInsertSql := "INSERT INTO novel_chapter_relation(novel_id, novel_chapters) VALUES (?,?)"
	_, err := NvlDB.Exec(NovelChapterRelationInsertSql, n.NovelID, n.NovelChapterToJson)
	if err != nil {
		fmt.Println("NvlDB.Exec NovelChapterRelationInsertSql Failed! err:", err)
		return
	}
}

func (n *NovelInfo) InsertIntoNovelClassifyRelation() {
	for _, v := range n.NovelClassifyToJson {
		ClassifyRelationInsertSql := "insert into novel_classify_relation(classify_id,novel_id) values(?,?)"
		_, err := NvlDB.Exec(ClassifyRelationInsertSql, v.NovelClassifyID, n.NovelID)
		if err != nil {
			fmt.Println("NvlDB.Exec ClassifyRelationInsertSql Failed! err:", err)
			return
		}
	}
}
func (n *NovelInfo) SelectNovelCategoryID() {
	SelectNovelCategoryIDSql := "select category_id from novel_category_info where category_name=?"
	err := NvlDB.QueryRow(SelectNovelCategoryIDSql, n.NovelCategoryName).Scan(&n.NovelCategoryID)
	if err != nil {
		fmt.Println("NvlDB.QueryRow SelectNovelCategoryIDSql Failed! err:", err)
		return
	}
}
func (n *NovelInfo) InsertIntoNovelCategoryRelation() {
	NovelCategoryRelationInsertSql := "insert into novel_category_relation(category_id,novel_id) values(?,?)"
	_, err := NvlDB.Exec(NovelCategoryRelationInsertSql, n.NovelCategoryID, n.NovelID)
	if err != nil {
		fmt.Println("NvlDB.Exec NovelCategoryRelationInsertSql Failed! err:", err)
		return
	}
}
