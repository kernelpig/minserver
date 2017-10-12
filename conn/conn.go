package conn

import (
	"fmt"
	"net"
	"sync"

	"wangqingang/server/handler"
	"wangqingang/server/pack"
)

var ConnList connList

func init() {
	ConnList.MapConns = make(map[net.Conn]*Conn, 0)
}

type connList struct {
	MapConns map[net.Conn]*Conn
	sync.Mutex
}

func (cl *connList) Get(c net.Conn) *Conn {
	if conn, ok := cl.MapConns[c]; ok {
		return conn
	}
	cl.MapConns[c] = &Conn{Conn: c, Queue: make([]string, 0)}
	return cl.MapConns[c]
}

type Conn struct {
	Queue      []string // 此连接上所有的待处理消息
	Conn       net.Conn // 连接
	sync.Mutex          // 消息队列锁
}

func (q *Conn) Proc() {
	for msg := q.First(); msg != ""; msg = q.First() {
		res := handler.MessageProcess(msg)
		resPack := []byte(pack.Pack(res))
		n, err := q.Conn.Write(resPack)
		if err != nil {
			// todo: add retry
			fmt.Println(err)
		} else if n != len(resPack) {
			// todo: write over check and proccess
			fmt.Println("write bytes: ", n)
		}
		fmt.Println("write over: ", res)
	}
}

func (q *Conn) First() string {
	q.Lock()
	defer q.Unlock()

	var msg string
	if len(q.Queue) <= 0 {
		return ""
	} else if len(q.Queue) == 1 {
		msg = q.Queue[0]
		q.Queue = make([]string, 0)
	} else {
		msg = q.Queue[0]
		q.Queue = q.Queue[1:]
	}
	fmt.Println("First(): ", msg)
	return msg
}

func (q *Conn) Append(msgs []string) {
	if msgs == nil || len(msgs) == 0 {
		return
	}
	q.Lock()
	defer q.Unlock()
	q.Queue = append(q.Queue, msgs...)
}
