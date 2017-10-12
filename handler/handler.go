package handler

import (
	"fmt"
	"strings"
)

type PutReq struct {
	ID      int
	Content string
	TTL     int64 // 生存时间单位秒
}

type GetReq struct {
	ID      int
	Refresh bool // 本次获取后是否重新刷新
}

type Response struct {
	Code        int
	Description string
}

const (
	actionPut = "put"
	actionGet = "get"
)

const (
	OK     = 0
	ErrXXX = 400
)

type handler func(string) string

var handlers map[string]handler

func init() {
	handlers = make(map[string]handler)
	handlers[actionPut] = PutHandler
	handlers[actionGet] = GetHandler
}

func PutHandler(request string) string {
	return ""
}

func GetHandler(request string) string {
	return ""
}

func DefaultHandler(request string) string {
	return ""
}

func MessageProcessor(message string) {
	action := strings.ToLower(message[:3])
	if handler, ok := handlers[action]; ok {
		response := handler(message[3:])
		fmt.Println(response)
	}
}
