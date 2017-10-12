package pack

import (
	"strings"

	"encoding/binary"
)

const (
	msgHeaderTypeFlag   = "www.cunxun.xin"
	msgHeaderLengthSize = 4
)

// 解包TLV格式
func UnPack(msgBytes []byte) ([]string, []byte) {
	cs := make([]string, 0)

	if msgBytes == nil || len(msgBytes) == 0 {
		return cs, msgBytes
	}
	msgHeaderIndex := strings.Index(string(msgBytes), msgHeaderTypeFlag)
	if msgHeaderIndex == -1 {
		return cs, msgBytes
	}

	msgContentLenIndex := msgHeaderIndex + len(msgHeaderTypeFlag)
	msgContentIndex := msgContentLenIndex + msgHeaderLengthSize

	if msgContentLenIndex+msgHeaderLengthSize > len(msgBytes) {
		return cs, msgBytes
	}
	msgContentLen := int(binary.BigEndian.Uint32(msgBytes[msgContentLenIndex : msgContentLenIndex+msgHeaderLengthSize]))
	if msgContentIndex+msgContentLen > len(msgBytes) {
		return cs, msgBytes
	} else if msgContentIndex+msgContentLen == len(msgBytes) {
		cs = append(cs, string(msgBytes[msgContentIndex:msgContentIndex+msgContentLen]))
		return cs, make([]byte, 0)
	}
	msgContent := msgBytes[msgContentIndex : msgContentIndex+msgContentLen]
	cs = append(cs, string(msgContent))
	subCs, subMsgBytes := UnPack(msgBytes[msgContentIndex+msgContentLen:])
	return append(cs, subCs...), subMsgBytes
}

// 封包TLV格式
func Pack(msgContent string) []byte {
	msgBuf := make([]byte, 0)
	msgContentLenBuf := make([]byte, 4)

	// 网络使用大端
	binary.BigEndian.PutUint32(msgContentLenBuf, uint32(len(msgContent)))

	msgBuf = append(msgBuf, []byte(msgHeaderTypeFlag)...)
	msgBuf = append(msgBuf, msgContentLenBuf...)
	msgBuf = append(msgBuf, []byte(msgContent)...)

	return msgBuf
}
