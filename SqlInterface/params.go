package main

import (
	"database/sql"
)

///home/datadrive/Gin
//漫画
const CartoonDsnTest = "root:123456@tcp(192.168.100.51:3306)/cartoon_test"
const CartoonDsnPro = "root:yxPvqJlbYBRrIs0z@tcp(110.92.66.88:6033)/res_cartoon_db"
const CartoonMD5Path = "./CartoonTitleMD5.txt"
const CartoonTitleSql = "select title,cartoon_id from cartoon_info"

var CartoonTitleMd5 map[string]string
var CtnDB *sql.DB

//小说
const NovelDsnTest = "root:123456@tcp(192.168.100.51:3306)/novel_test"
const NovelDsnPro = "root:yxPvqJlbYBRrIs0z@tcp(110.92.66.88:6033)/res_novel_db"
const NovelDsnTestPro = "root:ZsNice2020.@tcp(199.180.114.169:6033)/res_novel_db"
const NovelMD5Path = "./NovelTitleMD5.txt"
const NovelTitleSql = "select title,novel_id from novel_info Where novel_type=?"

var NovelTitleMd5 map[string]string
var NvlDB *sql.DB


//有声小说
const NovelSoundMD5Path = "./NovelSoundTitleMD5.txt"
const NovelSoundTitleSql = "select title,novel_id from novel_info Where novel_type=?"

var NovelSoundTitleMd5 map[string]string



