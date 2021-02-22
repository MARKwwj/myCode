package main

import (
	"fmt"
	"net"
)

func main() {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 30000})
	if err != nil {
		fmt.Println("net listenUDP failed,err:", err)
		return
	}
	defer udpConn.Close()
	for {
		var msgByte [1024]byte
		n, udpAddr, err := udpConn.ReadFromUDP(msgByte[:])
		if err != nil {
			fmt.Println("udpConn.ReadFromUDP(msgByte) failed,err:", err)
			return
		}
		fmt.Println("收到了:", string(msgByte[:n]))

		_, err = udpConn.WriteToUDP(msgByte[:n], udpAddr)
		fmt.Println("回复了:", string(msgByte[:n]))
	}
}
