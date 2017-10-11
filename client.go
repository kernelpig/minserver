package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	n, err := conn.Write([]byte("ABC"))
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}
