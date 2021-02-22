package main

import (
	"bufio"
	"fmt"
	pro "github.com/daily_exercise/tcp/protocol"
	"net"
)

func MsgAlter(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	//var buf [128]byte
	for {
		contentStr, err := pro.Decode(reader)
		if err != nil {
			fmt.Println("pro Decode failed,err:", err)
			return
		}
		fmt.Printf("收到消息:%v \n", contentStr)
		//
		//n, err := reader.Read(buf[:])
		//if err != nil {
		//	fmt.Println("reader read buf[:] failed,err:", err)
		//	break
		//}
		//fmt.Printf("收到消息:%v \n", string(buf[:n]))
	}
}
func main() {
	//建立监听
	listen, err := net.Listen("tcp", "127.0.0.1:8833")
	if err != nil {
		fmt.Println("net.Listen failed,err:", err)
		return
	}
	for {
		//等待连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen Accept failed,err", err)
			return
		}
		//接收消息
		go MsgAlter(conn)
	}
}
