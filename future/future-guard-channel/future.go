package main

import (
	"time"
	"errors"
)

type Process func() (interface{}, error)

func New(inFunc Process) Future {
	return newInner(make(chan struct{}), make(chan struct{}), inFunc)
}

func newInner(cancel chan struct{}, guard chan struct{}, inFunc Process) Future {
	f := futureImpl{
		done:   make(chan struct{}),
		cancel: cancel,
		guard: guard,
	}
	go func() {
		select {
		case <-f.done:
		//work was done, just exit
		case <-f.guard:
			close(f.cancel)
		}
	}()
	go func() {
		f.val, f.err = inFunc()
		close(f.done)
	}()
	return &f
}

type Step func(interface{}) (interface{}, error)

type Future interface {
	Get() (interface{}, error)

	GetUntil(d time.Duration) (interface{}, error)

	Then(Step) Future

	Cancel()

	IsCancelled() bool
}

type futureImpl struct {
	done   chan struct{}
	cancel chan struct{}
	guard  chan struct{}
	val    interface{}
	err    error
}

func (f *futureImpl) Get() (interface{}, error) {
	select {
	case <-f.done:
		return f.val, f.err
	case <-f.cancel:
	//on cancel, just fall out
	}
	return nil, nil
}

var FUTURE_TIMEOUT = errors.New("Your request has timed out")

func (f *futureImpl) GetUntil(d time.Duration) (interface{}, error) {
	select {
	case <-f.done:
		val, err := f.Get()
		return val, err
	case <-time.After(d):
		return nil, FUTURE_TIMEOUT
	case <-f.cancel:
	//on cancel, just fall out
	}
	return nil, nil
}

func (f *futureImpl) Then(next Step) Future {
	nextFuture := newInner(f.cancel, f.guard, func() (interface{}, error) {
		result, err := f.Get()
		if f.IsCancelled() || err != nil {
			return result, err
		}
		return next(result)
	})
	return nextFuture
}

func (f *futureImpl) Cancel() {
	select {
	case <-f.done:
	//already finished
	case <-f.cancel:
	//already cancelled
	default:
	//making request to cancel
	// this will serialize requests from multiple goroutines
			select {
			case <-f.done:
			case <-f.cancel:
			case f.guard <- struct{}{}:
			}
	}
}

func (f *futureImpl) IsCancelled() bool {
	select {
	case <-f.cancel:
		return true
	default:
		return false
	}
}

