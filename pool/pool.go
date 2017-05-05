package pool

import (
	"github.com/jonbodner/channel-talk-2/queue"
)

type Factory func() interface{}

type Pool interface {
	Borrow() interface{}
	Return(interface{})
}

type poolInner struct {
	items queue.Queue
}

func NewPool(f Factory, count int) Pool {
	pI := &poolInner{items: queue.MakeInfiniteQueue()}
	for i := 0; i < count; i++ {
		pI.items.Put(f())
	}
	return pI
}

func (pi *poolInner) Borrow() interface{} {
	item, _ := pi.items.Get()
	return item
}

func (pi *poolInner) Return(in interface{}) {
	pi.items.Put(in)
}
