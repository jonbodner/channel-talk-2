package queue

type Queue interface {
	Put(interface{})
	Get() (interface{}, bool)
	Close()
}

type queueInner struct {
	in  chan <- interface{}
	out <-chan interface{}
}

func (qi *queueInner) Get() (interface{}, bool) {
	v, ok := <-qi.out
	return v, ok
}

func (qi *queueInner) Put(val interface{}) {
	qi.in <- val
}

func (qi *queueInner) Close() {
	close(qi.in)
}
