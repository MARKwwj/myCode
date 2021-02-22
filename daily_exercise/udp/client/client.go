package main

import (
	"fmt"
	"net"
)

func main() {
	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 30000,
	})
	if err != nil {
		fmt.Println("net.DialUDP failed,err:", err)
		return
	}
	defer udpConn.Close()
	msgStr := "哈哈哈哈哈哈哈哈"
	_, err = udpConn.Write([]byte(msgStr))
	if err != nil {
		fmt.Println("udpConn.Write failed,err", err)
		return
	}
	fmt.Println("发送了:", msgStr)

	var data [1024]byte
	n, remoteAddr, err := udpConn.ReadFromUDP(data[:])
	if err != nil {
		fmt.Println("udpConn.ReadFromUDP failed,err:", err)
		return
	}
	fmt.Printf("recv:%v addr:%v count:%v\n", string(data[:n]), remoteAddr, n)

}
