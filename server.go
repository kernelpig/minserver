package main

import (
	"fmt"
	"net"
	"strings"
)

// 判断是否消息完成
func isMsgOver(bytes []byte) bool {
	if strings.Contains(string(bytes), "\n") {
		return true
	}
	return false
}

func connectionHandler(conn net.Conn) {
	buffer := make([]byte, 1)
	msgBuffer := make([]byte, 0)
	for {
		count, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("%+v", err)
			break
		}
		msgBuffer = append(msgBuffer, buffer...)
		if isMsgOver(buffer) {
			fmt.Println(msgBuffer)
			fmt.Println(string(msgBuffer))
			msgBuffer = make([]byte, 0)
		} else {
			msgBuffer = append(msgBuffer, buffer[:count]...)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("waiting for connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept connection failed: ", err)
			continue
		}
		connectionHandler(conn)
	}
}
