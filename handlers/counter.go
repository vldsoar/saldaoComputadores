package handlers

import "sync/atomic"

type Counter uint32

func (c *Counter) GetAndIncrement() uint32 {
	var next uint32

	for {
		next = uint32(*c) + 1
		if atomic.CompareAndSwapUint32((*uint32)(c), uint32(*c), next) {
			return next - 1
		}
	}
}

func (c *Counter) Get() uint32 {
	return atomic.LoadUint32((*uint32)(c))
}
