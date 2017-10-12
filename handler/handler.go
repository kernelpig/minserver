package handler

import (
	"fmt"

	"encoding/json"

	"wangqingang/server/cache"
	"wangqingang/server/proto"
)

type handler func(string) string

var handlers map[string]handler

func init() {
	handlers = make(map[string]handler)
	handlers[proto.ActionPut] = PutHandler
	handlers[proto.ActionGet] = GetHandler
}

func mustJson(object interface{}) string {
	bytes, err := json.Marshal(object)
	if err != nil {
		return fmt.Sprintf(`{"code": %d}`, proto.ErrJsonUnmarshal)
	}
	return string(bytes)
}

func PutHandler(request string) string {
	var req proto.PutReq
	var res proto.NormalRes
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		res.Code = proto.ErrJsonUnmarshal
		return mustJson(res)
	}
	res.Code = proto.OK
	return mustJson(res)
}

func GetHandler(request string) string {
	var req proto.GetReq
	var res proto.GetRes
	if err := json.Unmarshal([]byte(request), &req); err != nil {
		res.Code = proto.ErrJsonUnmarshal
		return mustJson(res)
	}
	res.Model = cache.Store.Get(req.ID, false)
	return mustJson(res)
}

func DefaultHandler(request string) string {
	return mustJson(proto.NormalRes{Code: proto.ErrActionNotSupport})
}

func MessageProcess(req string) string {
	action, content := proto.Unmarshal(req)
	if handler, ok := handlers[action]; ok {
		return handler(content)
	}
	return DefaultHandler(req)
}
