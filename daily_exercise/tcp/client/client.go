package main

import (
	"fmt"
	pro "github.com/daily_exercise/tcp/protocol"
	"net"
)

func main() {
	//拨号连接
	conn, err := net.Dial("tcp", "127.0.0.1:8833")
	if err != nil {
		fmt.Println("net Dial failed,err:", err)
		return
	}
	//发送消息
	defer conn.Close()
	//reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 10; i++ {
		//fmt.Println("请输入：")
		//var contentStr string
		//contentStr, err = reader.ReadString('\n')
		//if err != nil {
		//	fmt.Println("reader ReadString failed,err", err)
		//	break
		//}
		contentStr := "哈哈哈哈哈哈哈哈哈哈哈哈哈哈哈"
		var ContentByte []byte
		ContentByte, err = pro.Encode(contentStr)
		_, err = conn.Write(ContentByte)
		if err != nil {
			fmt.Println("conn write failed,err:", err)
			break
		}
	}
}
