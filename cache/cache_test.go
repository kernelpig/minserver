package cache

import (
	"sort"
	"strconv"
	"testing"
	"time"

	"math/rand"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func genTestData() []int {
	ids := make([]int, 0)
	timeout := time.Now().Add(-1 * time.Duration(1) * time.Hour)
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		n := rand.Intn(60)
		ttl := time.Duration(n) * time.Second
		model := Model{ID: n, TTL: ttl, CreateTime: timeout, Content: strconv.Itoa(n)}
		Store.Put(&model)
		ids = append(ids, model.ID)
	}
	sort.Ints(ids)
	return ids
}

// 定时器堆排序
func TestStore_Put(t *testing.T) {
	as := assert.New(t)
	Init()

	genIds := genTestData()
	popIds := make([]int, 0)

	popModels := Store.Flush()
	for _, m := range popModels {
		popIds = append(popIds, m.ID)
	}

	as.Equal(genIds, popIds)
}

// 获取数据，更新/不更新
func TestStore_Get(t *testing.T) {
	as := assert.New(t)
	Init()

	genIds := genTestData()

	m1 := Store.Get(genIds[0], false)
	as.Equal(1, m1.RefreshCount)

	m1 = Store.Get(genIds[0], true)
	as.Equal(2, m1.RefreshCount)
}

func patchTimeSleep(d time.Duration) {
	return
}

func TestStore_Tick(t *testing.T) {
	Init()

	genTestData()

	monkey.Patch(time.Sleep, patchTimeSleep)
	Store.Tick(true)
	monkey.UnpatchAll()
}
