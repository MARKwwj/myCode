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
	defer CtnDB.Close()
	var CtnInfo CartoonInfo
	if err := c.BindJSON(&CtnInfo); err != nil {
		fmt.Println(err)
	}
	// 先验证是否重复
	ok, err := PreventRepeating(CtnInfo.CartoonTitle, CartoonTitleMd5)
	if err != nil {
		fmt.Println("PreventRepeating Failed! err:", err)
		return
	}
	if CtnInfo.CartoonOnlyQuery == true {
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
		CtnInfo.InsertIntoCartoonClassifyInfo()
		CtnInfo.SelectCartoonCategoryID()
		CtnInfo.InsertIntoCartoonInfo()
		CtnInfo.InsertIntoCartoonChapterInfos()
		CtnInfo.InsertIntoCartoonChapterRelation()
		CtnInfo.InsertIntoCartoonClassifyRelation()
		CtnInfo.InsertIntoCartoonCategoryRelation()
		CtnInfo.InsertIntoCartoonPicInfo()
		NewTitleWriteToMap(CtnInfo.CartoonTitle, CartoonTitleMd5)
		c.JSON(http.StatusOK, gin.H{"CartoonID": CtnInfo.CartoonID})
	}
}

func (c *CartoonInfo) InsertIntoCartoonClassifyInfo() {
	var NvlClassify CartoonClassify
	var TagsStr string
	if c.CartoonClassifyNames == "" {
		c.CartoonClassifyToString = "[]"
		return
	}
	tagSlice := strings.Split(c.CartoonClassifyNames, ",")
	for _, ClassifyName := range tagSlice {
		var ClassifyID int64
		ClassifySql := "select classify_id from cartoon_classify_info where classify_name=?"
		err := CtnDB.QueryRow(ClassifySql, ClassifyName).Scan(&ClassifyID)
		switch {
		case err == sql.ErrNoRows:
			ClassifyInsertSql := "insert into cartoon_classify_info(category_id,classify_name) values(1,?)"
			result, err := CtnDB.Exec(ClassifyInsertSql, ClassifyName)
			if err != nil {
				fmt.Println("CtnDB Exec failed! err:", err)
				return
			}
			ClassifyID, _ = result.LastInsertId()
		case err != nil:
			fmt.Println("CtnDB queryRow failed! err:", err)
			return
		}
		//转 json
		NvlClassify.CartoonClassifyID = ClassifyID
		NvlClassify.CartoonClassifyName = ClassifyName
		c.CartoonClassifyToJson = append(c.CartoonClassifyToJson, NvlClassify)
	}
	tagsJson, _ := json.Marshal(c.CartoonClassifyToJson)
	TagsStr = string(tagsJson)
	c.CartoonClassifyToString = TagsStr
}

func (c *CartoonInfo) InsertIntoCartoonInfo() {
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
		"cartoon_classifies " +
		") VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	result, err := CtnDB.Exec(CartoonSql,
		c.CartoonTitle, c.CartoonAuthor, c.CartoonIntroduce,
		c.CartoonChapterTotalCount, c.CartoonChapterFreeCount, c.CartoonCoinPrice,
		c.CartoonStatus, c.CartoonPayType, c.CartoonEnable,
		c.CartoonPopularity, c.CartoonScore, c.CartoonFavoriteCount,
		c.CartoonCreateTime, c.CartoonCreator, c.CartoonClassifyToString,
	)
	if err != nil {
		fmt.Println("DB Exec Failed! err:", err)
		return
	}
	c.CartoonID, err = result.LastInsertId()
	if err != nil {
		fmt.Println("result LastInsertId Failed! err:", err)
		return
	}
}

func (c *CartoonInfo) InsertIntoCartoonChapterInfos() {
	for _, v := range c.CartoonChapterInfos {
		for x, _ := range v.CartoonPicInfos {
			v.CartoonChapterPicsSlice = append(v.CartoonChapterPicsSlice, x+1)
		}
		v.CartoonPicCount = len(v.CartoonChapterPicsSlice)
		byteSlice, _ := json.Marshal(v.CartoonChapterPicsSlice)
		v.CartoonChapterPicsString = string(byteSlice)
		CartoonChapterInfoInsertSql := "INSERT INTO cartoon_chapter_info(cartoon_id, chapter_title, chapter_pics,pic_count,sort ,chapter_byte_size) VALUES (?,?,?,?,?,?)"
		_, err := CtnDB.Exec(CartoonChapterInfoInsertSql, c.CartoonID, v.CartoonChapterName, v.CartoonChapterPicsString, v.CartoonPicCount, v.CartoonChapterSort, v.CartoonChapterByteSize)
		if err != nil {
			fmt.Println("CtnDB.Exec CartoonChapterInfoInsertSql Failed ! err:", err)
			return
		}
	}
}

func (c *CartoonInfo) InsertIntoCartoonChapterRelation() {
	CartoonChapterRelationInsertSql := "INSERT INTO cartoon_chapter_relation(cartoon_id, cartoon_chapters) VALUES (?,?)"
	_, err := CtnDB.Exec(CartoonChapterRelationInsertSql, c.CartoonID, c.CartoonChapterToJson)
	if err != nil {
		fmt.Println("CtnDB.Exec CartoonChapterRelationInsertSql Failed! err:", err)
		return
	}
}

func (c *CartoonInfo) InsertIntoCartoonClassifyRelation() {
	for _, v := range c.CartoonClassifyToJson {
		ClassifyRelationInsertSql := "insert into cartoon_classify_relation(classify_id,cartoon_id) values(?,?)"
		_, err := CtnDB.Exec(ClassifyRelationInsertSql, v.CartoonClassifyID, c.CartoonID)
		if err != nil {
			fmt.Println("CtnDB.Exec ClassifyRelationInsertSql Failed! err:", err)
			return
		}
	}
}

func (c *CartoonInfo) SelectCartoonCategoryID() {
	SelectCartoonCategoryIDSql := "select category_id from cartoon_category_info where category_name=?"
	err := CtnDB.QueryRow(SelectCartoonCategoryIDSql, c.CartoonCategoryName).Scan(&c.CartoonCategoryID)
	if err != nil {
		fmt.Println("CtnDB.QueryRow SelectCartoonCategoryIDSql Failed! err:", err)
		return
	}
}

func (c *CartoonInfo) InsertIntoCartoonCategoryRelation() {
	CartoonCategoryRelationInsertSql := "insert into cartoon_category_relation(category_id,cartoon_id) values(?,?)"
	_, err := CtnDB.Exec(CartoonCategoryRelationInsertSql, c.CartoonCategoryID, c.CartoonID)
	if err != nil {
		fmt.Println("CtnDB.Exec CartoonCategoryRelationInsertSql Failed! err:", err)
		return
	}
}

func (c *CartoonInfo) InsertIntoCartoonPicInfo() {
	for _, v := range c.CartoonChapterInfos {
		for _, y := range v.CartoonPicInfos {
			CartoonPicInfoInsertSql := "INSERT INTO cartoon_pic_info(chapter_id, pic_num, height, width) VALUES ( ?, ?, ?, ?)"
			_, err := CtnDB.Exec(CartoonPicInfoInsertSql, y.CartoonChapterID, y.CartoonPicNum, y.CartoonHeight, y.CartoonWidth)
			if err != nil {
				fmt.Println("CtnDB.Exec CartoonPicInfoInsertSql Failed! err:", err)
				return
			}
		}
	}
}
