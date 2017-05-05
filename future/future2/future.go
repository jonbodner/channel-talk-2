package main

type Process func() (interface{}, error)

func New(inFunc Process) Future {
	f := &futureImpl {
		done: make(chan struct{}),
	}
	go func() {
		f.val, f.err = inFunc()
		close(f.done)
	}()
	return f
}

type Future interface {
	Get() (interface{}, error)
}

type futureImpl struct {
	done chan struct{}
	val interface{}
	err error
}

func (f *futureImpl) Get() (interface{}, error) {
	<- f.done
	return f.val, f.err
}

