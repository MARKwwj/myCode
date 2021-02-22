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
	var novelID int
	if NvlInfo.NovelType == 1 {
		ok, novelID = PreventRepeating(NvlInfo.NovelTitle, NovelTitleMd5)
	} else if NvlInfo.NovelType == 2 {
		ok, novelID = PreventRepeating(NvlInfo.NovelTitle, NovelSoundTitleMd5)
	}
	if NvlInfo.NovelOnlyQuery == true {
		// true 重复 false 未重复
		if ok {
			c.JSON(http.StatusOK, gin.H{"repeat": true, "novelID": novelID})
			fmt.Printf("NovelTitle:%v,NovelType:%v ", NvlInfo.NovelTitle, NvlInfo.NovelType)
			fmt.Println("repeat: true")
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"repeat": false})
			fmt.Printf("NovelTitle:%v,NovelType:%v ", NvlInfo.NovelTitle, NvlInfo.NovelType)
			fmt.Println("repeat: false")
			return
		}
	} else {
		if NvlInfo.NovelTitle == "" || NvlInfo.NovelChapterTotalCount == 0 {
			return
		}
		if NvlInfo.NovelChapter == nil || NvlInfo.NovelByteSize == 0 {
			return
		}
		if ok {
			c.JSON(http.StatusOK, gin.H{"repeat": true, "novelID": novelID})
			fmt.Printf("NovelTitle:%v,NovelType:%v ", NvlInfo.NovelTitle, NvlInfo.NovelType)
			fmt.Println("repeat: true")
			return
		}
		//创建事务
		var err error
		NvlInfo.NvlDBTx, err = NvlDB.Begin()
		if err != nil {
			if NvlInfo.NvlDBTx != nil {
				NvlInfo.NvlDBTx.Rollback() //回滚
			}
			fmt.Println("NvlDB.Begin() failed,err:", err)
			return
		}
		bool1 := NvlInfo.SelectNovelCategoryID()
		bool2 := NvlInfo.InsertIntoClassifyInfo()
		bool3 := NvlInfo.InsertIntoNovelInfo()
		bool4 := NvlInfo.InsertIntoNovelChapterInfos()
		bool5 := NvlInfo.InsertIntoNovelChapterRelation()
		bool6 := NvlInfo.InsertIntoNovelClassifyRelation()
		bool7 := NvlInfo.InsertIntoNovelCategoryRelation()
		if bool1 && bool2 && bool3 && bool4 && bool5 && bool6 && bool7 {
			fmt.Println("事务提交啦...")
			NvlInfo.NvlDBTx.Commit()
			if NvlInfo.NovelType == 1 {
				NewTitleWriteToMap(NvlInfo.NovelTitle, NovelTitleMd5, NvlInfo.NovelID)
			} else if NvlInfo.NovelType == 2 {
				NewTitleWriteToMap(NvlInfo.NovelTitle, NovelSoundTitleMd5, NvlInfo.NovelID)
			}
			c.JSON(http.StatusOK, gin.H{"novelID": NvlInfo.NovelID, "chapterID": NvlInfo.NovelChapterToSlice})
		} else {
			NvlInfo.NvlDBTx.Rollback()
			fmt.Println("事务回滚啦...")
		}
	}
}
func (n *NovelInfo) InsertIntoClassifyInfo() bool {
	var NvlClassify NovelClassify
	var TagsStr string
	if n.NovelClassifyNames == "" {
		n.NovelClassifyToString = "[]"
		return true
	}
	tagSlice := strings.Split(n.NovelClassifyNames, ",")
	for _, ClassifyName := range tagSlice {
		var ClassifyID int64
		ClassifySql := "select classify_id from novel_classify_info where classify_name=? and category_id=?"
		err := NvlDB.QueryRow(ClassifySql, ClassifyName, n.NovelCategoryID).Scan(&ClassifyID)
		switch {
		case err == sql.ErrNoRows:
			ClassifyInsertSql := "insert into novel_classify_info(category_id,classify_name) values(?,?)"
			var result sql.Result
			result, err = n.NvlDBTx.Exec(ClassifyInsertSql, n.NovelCategoryID, ClassifyName)
			if err != nil {
				n.NvlDBTx.Rollback()
				fmt.Println("NvlDB Exec failed! err:", err)
				return false
			}

			if n.WhetherSuccess(result) == false {
				return false
			}

			ClassifyID, _ = result.LastInsertId()
		case err != nil:
			fmt.Println("NvlDB queryRow failed! err:", err)
			return false
		}
		//转 json
		NvlClassify.NovelCategoryId = n.NovelCategoryID
		NvlClassify.NovelClassifyID = ClassifyID
		NvlClassify.NovelClassifyName = ClassifyName
		n.NovelClassifyToJson = append(n.NovelClassifyToJson, NvlClassify)
	}
	tagsJson, _ := json.Marshal(n.NovelClassifyToJson)
	TagsStr = string(tagsJson)
	n.NovelClassifyToString = TagsStr
	return true
}
func (n *NovelInfo) InsertIntoNovelInfo() bool {
	n.NovelIntroduce = strings.TrimSpace(n.NovelIntroduce)
	if strings.Contains(n.NovelIntroduce,"简介："){
		n.NovelIntroduce = (strings.Split(n.NovelIntroduce, "简介："))[1]
	}
	n.NovelTitle = strings.TrimSpace(n.NovelTitle)
	n.NovelAuthor = strings.TrimSpace(n.NovelAuthor)

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
	n.NovelUpdateTime = time.Now().Format("2006-01-02 15:04:05")
	n.NovelFavoriteCount = rand.Intn(970) + 129
	n.NovelCreator = "reptile"
	n.NovelMachineID = 20000
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
		"novel_byte_size," +
		"favorite_count," +
		"novel_machine_id," +
		"chapter_update_time) " +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := n.NvlDBTx.Exec(NovelInfoInsertSql,
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
		n.NovelByteSize,
		n.NovelFavoriteCount,
		n.NovelMachineID,
		n.NovelUpdateTime,
	)
	if err != nil {
		n.NvlDBTx.Rollback()
		fmt.Println("NvlDB.Exec NovelInfoInsertSql Failed! err:", err)
		return false
	}

	if n.WhetherSuccess(result) == false {
		return false
	}

	var novelID int64
	novelID, err = result.LastInsertId()
	n.NovelID = int(novelID)
	if err != nil {
		fmt.Println("NovelID  result.LastInsertId() Failed! err:", err)
		return false
	}
	return true

}

func (n *NovelInfo) InsertIntoNovelChapterInfos() bool {
	for _, v := range n.NovelChapter {
		v.NovelChapterName = strings.TrimSpace(v.NovelChapterName)
		NovelChapterInfoInsertSql := "INSERT INTO novel_chapter_info(novel_id, chapter_name, chapter_no, chapter_byte_size,chapter_time) VALUES (?,?,?,?,?)"
		result, err := n.NvlDBTx.Exec(NovelChapterInfoInsertSql, n.NovelID, v.NovelChapterName, v.NovelChapterSort, v.NovelChapterByteSize, v.NovelChapterTime)
		if err != nil {
			n.NvlDBTx.Rollback()
			fmt.Println("NvlDB.Exec NovelChapterInfoInsertSql Failed ! err:", err)
			return false
		}

		if n.WhetherSuccess(result) == false {
			return false
		}

		chapterID, _ := result.LastInsertId()
		n.NovelChapterToSlice = append(n.NovelChapterToSlice, int(chapterID))
	}
	byteSlice, _ := json.Marshal(n.NovelChapterToSlice)
	n.NovelChapterToJson = string(byteSlice)
	return true
}

func (n *NovelInfo) InsertIntoNovelChapterRelation() bool {
	NovelChapterRelationInsertSql := "INSERT INTO novel_chapter_relation(novel_id, novel_chapters) VALUES (?,?)"
	result, err := n.NvlDBTx.Exec(NovelChapterRelationInsertSql, n.NovelID, n.NovelChapterToJson)
	if err != nil {
		n.NvlDBTx.Rollback()
		fmt.Println("NvlDB.Exec NovelChapterRelationInsertSql Failed! err:", err)
		return false
	}
	return n.WhetherSuccess(result)
}

func (n *NovelInfo) InsertIntoNovelClassifyRelation() bool {
	for _, v := range n.NovelClassifyToJson {
		ClassifyRelationInsertSql := "insert into novel_classify_relation(classify_id,novel_id) values(?,?)"
		result, err := n.NvlDBTx.Exec(ClassifyRelationInsertSql, v.NovelClassifyID, n.NovelID)
		if err != nil {
			n.NvlDBTx.Rollback()
			fmt.Println("NvlDB.Exec ClassifyRelationInsertSql Failed! err:", err)
			return false
		}
		if n.WhetherSuccess(result) == false {
			return false
		}
	}
	return true
}
func (n *NovelInfo) SelectNovelCategoryID() bool {
	SelectNovelCategoryIDSql := "select category_id from novel_category_info where category_name=?"
	err := NvlDB.QueryRow(SelectNovelCategoryIDSql, n.NovelCategoryName).Scan(&n.NovelCategoryID)
	if err != nil {
		fmt.Println("NvlDB.QueryRow SelectNovelCategoryIDSql Failed! err:", err)
		return false
	}
	if n.NovelCategoryID == 0 {
		return false
	} else {
		return true
	}
}
func (n *NovelInfo) InsertIntoNovelCategoryRelation() bool {
	NovelCategoryRelationInsertSql := "insert into novel_category_relation(category_id,novel_id) values(?,?)"
	result, err := n.NvlDBTx.Exec(NovelCategoryRelationInsertSql, n.NovelCategoryID, n.NovelID)
	if err != nil {
		n.NvlDBTx.Rollback()
		fmt.Println("NvlDB.Exec NovelCategoryRelationInsertSql Failed! err:", err)
		return false
	}
	return n.WhetherSuccess(result)
}
func (n *NovelInfo) WhetherSuccess(result sql.Result) bool {
	affRow, err := result.RowsAffected()
	if err != nil {
		n.NvlDBTx.Rollback()
		fmt.Println("ClassifyRelationInsertSql result.RowsAffected() Failed! err:", err)
		return false
	}
	if affRow != 1 {
		return false
	}
	return true
}
