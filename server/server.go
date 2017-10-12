package main

import (
	"fmt"
	"net"

	"wangqingang/server/conn"
	"wangqingang/server/pack"
)

func connectionHandler(netConn net.Conn) {

	c := conn.ConnList.Get(netConn)

	// todo: buffer设置大点，此处设置小方便测试
	tmpBuffer := make([]byte, 1)
	var msgBuffer []byte
	var msgs []string
	for {
		count, err := netConn.Read(tmpBuffer)
		if err != nil {
			fmt.Println(err)
			// todo: 生产环境不break
			break
		}
		tmpBuffer = tmpBuffer[:count]
		msgBuffer = append(msgBuffer, tmpBuffer...)
		msgs, msgBuffer = pack.UnPack(msgBuffer)
		if len(msgs) != 0 {
			// todo: 异步处理
			c.Append(msgs)
			c.Proc()
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
