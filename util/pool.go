package util

import (
	"sync/atomic"
)

type IntPool struct {
	values chan int
	size   int
	idle   int32
}

func NewIntPool(size int) *IntPool {
	pool := &IntPool{
		values: make(chan int, size),
		size:   size,
		idle:   int32(size),
	}

	pool.init()
	return pool
}

func (p *IntPool) init() {
	for i := 0; i < cap(p.values); i++ {
		p.values <- i
	}
}

func (p *IntPool) Get() int {
	val := <-p.values
	atomic.AddInt32(&p.idle, int32(-1))
	return val
}

func (p *IntPool) Put(obj int) {
	p.values <- obj
	atomic.AddInt32(&p.idle, int32(1))
}
