package cache

import (
	"fmt"
	"sync"
	"time"

	"container/heap"
)

var Store store

func Init() {
	Store.mapBuffer = make(map[int]*Model)
	Store.timer = make([]*Model, 0)
	heap.Init(&Store.timer)
}

// 数据存储: 最小堆+map，最小堆负责老化， map负责存储数据
type Model struct {
	timerIndex   int // 本元素在heap中的索引
	ID           int
	Content      string
	CreateTime   time.Time
	RefreshCount int
	TTL          time.Duration // 生存周期
}

func (m *Model) GetDeadline() time.Time {
	duration := time.Duration(int(m.TTL.Seconds())*m.RefreshCount) * time.Second
	return m.CreateTime.Add(duration)
}

func (m *Model) IsTimeout() bool {
	return time.Now().After(m.GetDeadline())
}

type timer []*Model
type store struct {
	mapBuffer map[int]*Model
	timer     timer
	sync.Mutex
}

func (h timer) Len() int {
	return len(h)
}

func (h timer) Less(i, j int) bool {
	return h[i].GetDeadline().Before(h[j].GetDeadline())
}

func (h timer) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].timerIndex, h[j].timerIndex = j, i
}

func (h *timer) Push(x interface{}) {
	*h = append(*h, x.(*Model))
}

func (h *timer) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *timer) Add(x interface{}) {
	heap.Push(h, x)
	heap.Fix(h, h.Len()-1)
}

func (h *timer) String() {
	fmt.Println("/////////////// all data ////////////")
	for i := 0; i < h.Len(); i++ {
		x := (*h)[i]
		fmt.Println(x.ID, x.TTL, x.GetDeadline())
	}
}

func (c *store) Put(m *Model) {
	c.Lock()
	defer c.Unlock()

	m.RefreshCount = 1
	c.mapBuffer[m.ID] = m
	(&c.timer).Add(m)
}

func (c *store) Get(id int, refresh bool) *Model {
	c.Lock()
	defer c.Unlock()

	m, ok := c.mapBuffer[id]
	if !ok {
		return nil
	}
	if !refresh {
		return m
	}
	m.RefreshCount += 1
	return m
}

func (c *store) GetTop() *Model {
	if c.timer.Len() > 0 {
		return c.timer[0]
	}
	return nil
}

func (c *store) DelTop() {
	c.Lock()
	defer c.Unlock()

	c.timer.String()
	if m, ok := heap.Pop(&c.timer).(*Model); ok {
		fmt.Println("del top: ", m.ID, m.GetDeadline())
		delete(c.mapBuffer, m.ID)
	}
}

func (c *store) Flush() []*Model {
	result := make([]*Model, 0)
	for Store.timer.Len() > 0 {
		if m, ok := heap.Pop(&Store.timer).(*Model); ok {
			result = append(result, m)
		}
	}
	return result
}

func (c *store) Tick(stop bool) {
	for {
		m := c.GetTop()
		if m == nil {
			// 没有数据设置了退出则退出
			if stop {
				return
			}
		} else if m != nil && m.IsTimeout() {
			c.DelTop()
		}
		// 延时1秒， 精度为1秒
		time.Sleep(time.Duration(1) * time.Second)
	}
}
