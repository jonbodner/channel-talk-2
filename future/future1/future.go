package main

type Process func() (interface{}, error)

func New(inFunc Process) Future {
	//do magic stuff
	return &futureImpl{}
}

type Future interface {
	Get() (interface{}, error)
}

type futureImpl struct {
	//this will have stuff eventually
}

func (f *futureImpl) Get() (interface{}, error) {
	//this will have stuff eventually
	return nil, nil
}

