package pack

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	msgContent   = "I'm a test message package"
	msgExtraHead = "extra_data_head"
	msgExtraTail = "extra_data_tail"
	msgCount     = 5
)

// 单包
func TestPack(t *testing.T) {
	as := assert.New(t)

	content := msgContent
	flow := Pack(content)

	cs, flow := UnPack(flow)
	as.Equal(1, len(cs))
	as.Equal(content, cs[0])
	as.Equal(0, len(flow))
}

// 多包
func TestUnPack(t *testing.T) {
	as := assert.New(t)

	content := msgContent
	flow := Pack(content)

	for i := 0; i < msgCount; i++ {
		flow = append(flow, flow...)
	}

	cs, leftFlow := UnPack(flow)
	as.Equal(int(math.Pow(2, msgCount)), len(cs))
	as.Equal(content, cs[0])
	as.Equal(0, len(leftFlow))
}

// 冗余数据+单包
func TestUnPack2(t *testing.T) {
	as := assert.New(t)

	content := msgContent
	flow := Pack(content)

	cs, leftFlow := UnPack(append([]byte(msgExtraHead), flow...))
	as.Equal(1, len(cs))
	as.Equal(0, len(leftFlow))
}

// 冗余数据+单包+冗余数据
func TestUnPack3(t *testing.T) {
	as := assert.New(t)

	content := msgContent
	flow := Pack(content)
	flow = append([]byte(msgExtraHead), flow...)
	flow = append(flow, []byte(msgExtraTail)...)

	cs, leftFlow := UnPack(flow)
	as.Equal(1, len(cs))
	as.Equal([]byte(msgExtraTail), leftFlow)
}

// 冗余数据+N包+冗余数据
func TestUnPack4(t *testing.T) {
	as := assert.New(t)

	content := msgContent
	flow := Pack(content)
	for i := 0; i < msgCount; i++ {
		flow = append(flow, flow...)
	}

	flow = append([]byte(msgExtraHead), flow...)
	flow = append(flow, []byte(msgExtraTail)...)

	cs, leftFlow := UnPack(flow)
	as.Equal(int(math.Pow(2, msgCount)), len(cs))
	as.Equal([]byte(msgExtraTail), leftFlow)
}
