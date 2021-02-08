package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func InitDB(Dsn string) (DB *sql.DB) {
	//不会校验账号密码是否正确
	DB, err := sql.Open("mysql", Dsn)
	if err != nil {
		fmt.Println("sql open failed! err:", err)
		return nil
	}
	//尝试与数据库连接，并校验dsn是否正确
	err = DB.Ping()
	if err != nil {
		fmt.Println("DB ping failed! err:", err)
		return nil
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Second * 500)
	return DB
}
func StrTomMD5(s string) string {
	m := md5.New()
	m.Write([]byte (s))
	return hex.EncodeToString(m.Sum(nil))
}

//放重复校验
func TitleWriteToTxt(titleSql, fileName string, Type int, DB *sql.DB) {
	//titleSql := "select title from cartoon_info"
	titleMd5 := make(map[string]string)
	var rows *sql.Rows
	var err error
	if Type != 0 {
		rows, err = DB.Query(titleSql, Type)
	} else {
		rows, err = DB.Query(titleSql)
	}
	if err != nil {
		fmt.Println("DB Query failed! err:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var title string
		var MD5 string
		var novelID int
		_ = rows.Scan(&title, &novelID)
		novelIDStr := strconv.Itoa(novelID)
		MD5 = StrTomMD5(title)
		titleMd5[MD5] = title + "->" + novelIDStr
	}
	//写入文件
	fileObj, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("fileObj is failed! err:", err)
		return
	}
	//转string
	dataType, _ := json.Marshal(titleMd5)
	dataString := string(dataType)
	_, err = fileObj.WriteString(dataString)
	//fmt.Println("writeBefore:", dataString)
	if err != nil {
		fmt.Println("fileObj WriteString failed! err:", err)
		return
	}
	fileObj.Close()
}
func TxtReadToMap(MD5Path string, titleMd5 map[string]string) {
	var b []byte
	var err error
	b, err = ioutil.ReadFile(MD5Path)
	if err != nil {
		fmt.Println("ioutil ReadFIle failed! err:", err)
		return
	}
	err = json.Unmarshal(b, &titleMd5)
	if err != nil {
		fmt.Println("json unmarshall failed! err:", err)
		return
	}
}

func NewTitleWriteToMap(titleName string, titleMd5 map[string]string, ID int) {
	CurTitleMd5 := StrTomMD5(titleName)
	IDStr := strconv.Itoa(ID)
	//rlk.Lock()
	titleMd5[CurTitleMd5] = titleName + "->" + IDStr
	//rlk.Unlock()
}
func PreventRepeating(titleName string, titleMd5 map[string]string) (bool, int) {
	CurTitleMd5 := StrTomMD5(titleName)
	//rlk.RLock()
	value, ok := titleMd5[CurTitleMd5]
	//rlk.RUnlock()
	if ok {
		result := (strings.Split(value, "->"))[1]
		ID, err := strconv.Atoi(result)
		if err != nil {
			fmt.Println("strconv.Atoi(result) failed! err:", err)
			return ok, 0
		}
		return ok, ID
	}
	return ok, 0
}
