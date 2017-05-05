package main

import (
	"time"
	"errors"
)

type Process func() (interface{}, error)

func New(inFunc Process) Future {
	f := &futureImpl{
		done:   make(chan struct{}),
	}
	go func() {
		f.val, f.err = inFunc()
		close(f.done)
	}()
	return f
}

type Step func(interface{}) (interface{}, error)

type Future interface {
	Get() (interface{}, error)

	GetUntil(d time.Duration) (interface{}, error)

	Then(Step) Future
}

type futureImpl struct {
	done   chan struct{}
	val    interface{}
	err    error
}

func (f *futureImpl) Get() (interface{}, error) {
	<-f.done
	return f.val, f.err
}

var FUTURE_TIMEOUT = errors.New("Your request has timed out")

func (f *futureImpl) GetUntil(d time.Duration) (interface{}, error) {
	select {
	case <-f.done:
		val, err := f.Get()
		return val, err
	case <-time.After(d):
		return nil, FUTURE_TIMEOUT
	}
	// This should never be executed
	return nil, nil
}

func (f *futureImpl) Then(next Step) Future {
	nextFuture := New(func() (interface{}, error) {
		result, err := f.Get()
		if err != nil {
			return result, err
		}
		return next(result)
	})
	return nextFuture
}

