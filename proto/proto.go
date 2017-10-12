package proto

import (
	"fmt"
	"strings"

	"wangqingang/server/cache"
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

type NormalRes struct {
	Code int
	More string
}

type GetRes struct {
	Code  int
	Model *cache.Model
}

const (
	ActionPut = "put"
	ActionGet = "get"
)

const (
	OK                  = 0
	ErrJsonUnmarshal    = 400
	ErrActionNotSupport = 401
)

func Unmarshal(message string) (string, string) {
	action := strings.ToLower(message[:3])
	return action, message[3:]
}

func Marshal(action, content string) string {
	return fmt.Sprintf("%s%s", action[:3], content)
}
