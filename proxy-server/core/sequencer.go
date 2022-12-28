package core

import (
	"math"
	"sync/atomic"
)

type Sequencer struct {
	current uint32
}

// NewSequencer 创建TCP序列号
func NewSequencer() *Sequencer {
	return &Sequencer{current: math.MaxUint32}
}

// Next 返回下一个TCP序列号
func (t *Sequencer) Next() uint32 {
	value := atomic.AddUint32(&t.current, 1) // 增加
	return value
}
