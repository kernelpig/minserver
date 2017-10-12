package client

import (
	"encoding/json"
	"fmt"
	"net"

	"wangqingang/server/pack"
	"wangqingang/server/proto"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	req := proto.PutReq{
		ID:      1,
		Content: "put request from client",
		TTL:     30,
	}

	reqJson, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}

	msg := proto.Marshal(proto.ActionPut, string(reqJson))

	n, err := conn.Write([]byte(pack.Pack(string(msg))))
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}
