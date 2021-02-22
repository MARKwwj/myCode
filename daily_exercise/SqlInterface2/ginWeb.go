package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func main() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "success"})
	})
	//小说 漫画 接口
	InitTitle()
	r.POST("/CartoonInfos", CartoonMainFunc)
	r.POST("/NovelInfos", NovelMainFunc)
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8002")
}

func InitTitle() {
	//小说
	NvlDB = InitDB(NovelDsnTest)
	NovelTitleMd5 = make(map[string]string, 2000)
	TitleWriteToTxt(NovelTitleSql, NovelMD5Path, 1, NvlDB)
	TxtReadToMap(NovelMD5Path, NovelTitleMd5)
	//有声小说
	NovelSoundTitleMd5 = make(map[string]string, 2000)
	TitleWriteToTxt(NovelSoundTitleSql, NovelSoundMD5Path, 2, NvlDB)
	TxtReadToMap(NovelSoundMD5Path, NovelSoundTitleMd5)
	//漫画
	CtnDB = InitDB(CartoonDsnTest)
	CartoonTitleMd5 = make(map[string]string, 2000)
	TitleWriteToTxt(CartoonTitleSql, CartoonMD5Path, 0, CtnDB)
	TxtReadToMap(CartoonMD5Path, CartoonTitleMd5)
}
