package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const ResDirPath = "/home/datadrive/resources/res2/smallvideo/crypted/videos"
//const ResDirPath = "D:\\desktop\\r1"

func main() {
	fileInfos, err := ioutil.ReadDir(ResDirPath)
	if err != nil {
		fmt.Println("ioutil ReadDir failed! err:", err)
		return
	}
	for _, v := range fileInfos {
		CurDirPath := path.Join(ResDirPath, v.Name())
		fmt.Println(CurDirPath)
		DirNewName := (strings.Split(v.Name(), "-A"))[0]
		DirNewPath := path.Join(ResDirPath, DirNewName)
		err = os.Rename(CurDirPath, DirNewPath)
		if err != nil {
			fmt.Println("Rename failed! err:", err)
			continue
		}
	}
}
