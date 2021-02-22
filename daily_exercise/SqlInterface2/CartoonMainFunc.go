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

func CartoonMainFunc(c *gin.Context) {
	var CtnInfo CartoonInfo
	if err := c.BindJSON(&CtnInfo); err != nil {
		fmt.Println(err)
	}

	// 先验证是否重复
	//1、文字小说  2、有声小说
	var ok bool
	var CartoonID int
	ok, CartoonID = PreventRepeating(CtnInfo.CartoonTitle, CartoonTitleMd5)

	if CtnInfo.CartoonOnlyQuery == true {
		// true 重复 false 未重复
		if ok {
			c.JSON(http.StatusOK, gin.H{"repeat": true, "CartoonID": CartoonID})
			fmt.Printf("CartoonTitle:%v,", CtnInfo.CartoonTitle)
			fmt.Println("repeat: true")
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"repeat": false})
			fmt.Printf("CartoonTitle:%v, ", CtnInfo.CartoonTitle)
			fmt.Println("repeat: false")
			return
		}
	} else {
		if CtnInfo.CartoonTitle == "" || CtnInfo.CartoonChapterTotalCount == 0 || CtnInfo.CartoonChapterInfos == nil {
			return
		}

		if ok {
			c.JSON(http.StatusOK, gin.H{"repeat": true, "CartoonID": CartoonID})
			fmt.Printf("CartoonTitle:%v, ", CtnInfo.CartoonTitle)
			fmt.Println("repeat: true")
			return
		}

		//创建事务
		var err error
		CtnInfo.CtnDBTx, err = CtnDB.Begin()
		if err != nil {
			if CtnInfo.CtnDBTx != nil {
				CtnInfo.CtnDBTx.Rollback() //回滚
			}
			fmt.Println("CtnDB.Begin() failed,err:", err)
			return
		}

		bool1 := CtnInfo.SelectCartoonCategoryID()
		bool2 := CtnInfo.InsertIntoCartoonClassifyInfo()
		bool3 := CtnInfo.InsertIntoCartoonInfo()
		bool4 := CtnInfo.InsertIntoCartoonChapterInfos()
		bool5 := CtnInfo.InsertIntoCartoonChapterRelation()
		bool6 := CtnInfo.InsertIntoCartoonClassifyRelation()
		bool7 := CtnInfo.InsertIntoCartoonCategoryRelation()
		bool8 := CtnInfo.InsertIntoCartoonPicInfo()
		if bool1 && bool2 && bool3 && bool4 && bool5 && bool6 && bool7 && bool8 {
			CtnInfo.CtnDBTx.Commit()
			fmt.Println("事务提交啦...")
			NewTitleWriteToMap(CtnInfo.CartoonTitle, CartoonTitleMd5, CtnInfo.CartoonID)
			c.JSON(http.StatusOK, gin.H{"CartoonID": CtnInfo.CartoonID, "chapterID": CtnInfo.CartoonChapterToSlice})
		} else {
			CtnInfo.CtnDBTx.Rollback()
			fmt.Println("事务回滚啦...")
		}
	}
}

func (c *CartoonInfo) InsertIntoCartoonClassifyInfo() bool {
	var CtnClassify CartoonClassify
	var TagsStr string
	if c.CartoonClassifyNames == "" {
		c.CartoonClassifyToString = "[]"
		return true
	}
	tagSlice := strings.Split(c.CartoonClassifyNames, ",")
	for _, ClassifyName := range tagSlice {
		var ClassifyID int64
		ClassifySql := "select classify_id from cartoon_classify_info where classify_name=?"
		err := CtnDB.QueryRow(ClassifySql, ClassifyName).Scan(&ClassifyID)
		switch {
		case err == sql.ErrNoRows:
			ClassifyInsertSql := "insert into cartoon_classify_info(category_id,classify_name) values(1,?)"
			var result sql.Result
			result, err = c.CtnDBTx.Exec(ClassifyInsertSql, ClassifyName)
			if err != nil {
				fmt.Println("CtnDB Exec failed! err:", err)
				return false
			}
			if c.WhetherSuccess(result) == false {
				return false
			}

			ClassifyID, _ = result.LastInsertId()
		case err != nil:
			fmt.Println("CtnDB queryRow failed! err:", err)
			return false
		}
		//转 json
		CtnClassify.CartoonCategoryId = c.CartoonCategoryID
		CtnClassify.CartoonClassifyID = ClassifyID
		CtnClassify.CartoonClassifyName = ClassifyName
		c.CartoonClassifyToJson = append(c.CartoonClassifyToJson, CtnClassify)
	}
	tagsJson, _ := json.Marshal(c.CartoonClassifyToJson)
	TagsStr = string(tagsJson)
	c.CartoonClassifyToString = TagsStr
	return true
}

func (c *CartoonInfo) InsertIntoCartoonInfo() bool {
	c.CartoonIntroduce = strings.TrimSpace(c.CartoonIntroduce)
	if strings.Contains(c.CartoonIntroduce, "简介：") {
		c.CartoonIntroduce = (strings.Split(c.CartoonIntroduce, "简介："))[1]
	}
	c.CartoonTitle = strings.TrimSpace(c.CartoonTitle)
	c.CartoonAuthor = strings.TrimSpace(c.CartoonAuthor)

	if c.CartoonAuthor == "" {
		c.CartoonAuthor = "来自于网络"
	}
	if c.CartoonIntroduce == "" {
		c.CartoonIntroduce = c.CartoonTitle
	}
	c.CartoonCoinPrice = 4
	c.CartoonChapterFreeCount = 1
	c.CartoonPopularity = rand.Intn(87000) + 12900
	c.CartoonScore = (float32(rand.Intn(9)) / 10) + float32(8.8)
	c.CartoonFavoriteCount = rand.Intn(870) + 129
	c.CartoonEnable = 1
	c.CartoonPayType = 1
	c.CartoonCreator = "reptile"
	c.CartoonCreateTime = time.Now().Format("2006-01-02 15:04:05")
	//漫画所有章节的Id json数组
	for k, _ := range c.CartoonChapterInfos {
		c.CartoonChapterToSlice = append(c.CartoonChapterToSlice, k+1)
	}
	byteSlice, _ := json.Marshal(c.CartoonChapterToSlice)
	c.CartoonChapterToJson = string(byteSlice)
	CartoonSql := "INSERT INTO cartoon_info(" +
		"title, " +
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
		"favorite_count, " +
		"create_time, " +
		"cartoon_creator, " +
		"cartoon_classifies," +
		"last_update_time " +
		") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := CtnDB.Exec(CartoonSql,
		c.CartoonTitle, c.CartoonAuthor, c.CartoonIntroduce,
		c.CartoonChapterTotalCount, c.CartoonChapterFreeCount, c.CartoonCoinPrice,
		c.CartoonStatus, c.CartoonPayType, c.CartoonEnable,
		c.CartoonPopularity, c.CartoonScore, c.CartoonFavoriteCount,
		c.CartoonCreateTime, c.CartoonCreator, c.CartoonClassifyToString, c.CartoonCreateTime,
	)
	if err != nil {
		fmt.Println("DB Exec Failed! err:", err)
		return false
	}
	if c.WhetherSuccess(result) == false {
		return false
	}
	var CtnID int64
	CtnID, err = result.LastInsertId()
	c.CartoonID = int(CtnID)
	if err != nil {
		fmt.Println("result LastInsertId Failed! err:", err)
		return false
	}
	return true
}

func (c *CartoonInfo) InsertIntoCartoonChapterInfos() bool {
	for _, v := range c.CartoonChapterInfos {
		for x, _ := range v.CartoonPicInfos {
			v.CartoonChapterPicsSlice = append(v.CartoonChapterPicsSlice, x+1)
		}
		v.CartoonPicCount = len(v.CartoonChapterPicsSlice)
		byteSlice, _ := json.Marshal(v.CartoonChapterPicsSlice)
		v.CartoonChapterPicsString = string(byteSlice)
		CartoonChapterInfoInsertSql := "INSERT INTO cartoon_chapter_info(cartoon_id, chapter_title, chapter_pics,pic_count,sort ,chapter_byte_size) VALUES (?,?,?,?,?,?)"
		result, err := c.CtnDBTx.Exec(CartoonChapterInfoInsertSql, c.CartoonID, v.CartoonChapterName, v.CartoonChapterPicsString, v.CartoonPicCount, v.CartoonChapterSort, v.CartoonChapterByteSize)
		if err != nil {
			fmt.Println("CtnDB.Exec CartoonChapterInfoInsertSql Failed ! err:", err)
			return false
		}
		if c.WhetherSuccess(result) == false {
			return false
		}
	}
	return true
}

func (c *CartoonInfo) InsertIntoCartoonChapterRelation() bool {
	CartoonChapterRelationInsertSql := "INSERT INTO cartoon_chapter_relation(cartoon_id, cartoon_chapters) VALUES (?,?)"
	result, err := c.CtnDBTx.Exec(CartoonChapterRelationInsertSql, c.CartoonID, c.CartoonChapterToJson)
	if err != nil {
		fmt.Println("CtnDB.Exec CartoonChapterRelationInsertSql Failed! err:", err)
		return false
	}
	if c.WhetherSuccess(result) == false {
		return false
	}
	return true
}

func (c *CartoonInfo) InsertIntoCartoonClassifyRelation() bool {
	for _, v := range c.CartoonClassifyToJson {
		ClassifyRelationInsertSql := "insert into cartoon_classify_relation(classify_id,cartoon_id) values(?,?)"
		result, err := c.CtnDBTx.Exec(ClassifyRelationInsertSql, v.CartoonClassifyID, c.CartoonID)
		if err != nil {
			fmt.Println("CtnDB.Exec ClassifyRelationInsertSql Failed! err:", err)
			return false
		}
		if c.WhetherSuccess(result) == false {
			return false
		}
	}

	return true
}

func (c *CartoonInfo) SelectCartoonCategoryID() bool {
	SelectCartoonCategoryIDSql := "select category_id from cartoon_category_info where category_name=?"
	err := CtnDB.QueryRow(SelectCartoonCategoryIDSql, c.CartoonCategoryName).Scan(&c.CartoonCategoryID)
	if err != nil {
		fmt.Println("CtnDB.QueryRow SelectCartoonCategoryIDSql Failed! err:", err)
		return false
	}
	if c.CartoonCategoryID == 0 {
		return false
	} else {
		return true
	}
}

func (c *CartoonInfo) InsertIntoCartoonCategoryRelation() bool {
	CartoonCategoryRelationInsertSql := "insert into cartoon_category_relation(category_id,cartoon_id) values(?,?)"
	result, err := c.CtnDBTx.Exec(CartoonCategoryRelationInsertSql, c.CartoonCategoryID, c.CartoonID)
	if err != nil {
		fmt.Println("CtnDB.Exec CartoonCategoryRelationInsertSql Failed! err:", err)
		return false
	}
	if c.WhetherSuccess(result) == false {
		return false
	}
	return true
}

func (c *CartoonInfo) InsertIntoCartoonPicInfo() bool {
	for _, v := range c.CartoonChapterInfos {
		for _, y := range v.CartoonPicInfos {
			CartoonPicInfoInsertSql := "INSERT INTO cartoon_pic_info(chapter_id, pic_num, height, width) VALUES ( ?, ?, ?, ?)"
			result, err := c.CtnDBTx.Exec(CartoonPicInfoInsertSql, y.CartoonChapterID, y.CartoonPicNum, y.CartoonHeight, y.CartoonWidth)
			if err != nil {
				fmt.Println("CtnDB.Exec CartoonPicInfoInsertSql Failed! err:", err)
				return false
			}
			if c.WhetherSuccess(result) == false {
				return false
			}
		}
	}
	return true
}
func (c *CartoonInfo) WhetherSuccess(result sql.Result) bool {
	affRow, err := result.RowsAffected()
	if err != nil {
		c.CtnDBTx.Rollback()
		fmt.Println("ClassifyRelationInsertSql result.RowsAffected() Failed! err:", err)
		return false
	}
	if affRow != 1 {
		return false
	}
	return true
}
