package main

import (
	"fmt"
	"testing"
	"time"
	"sync"
	"sync/atomic"
	"runtime"
)

func doSomethingThatTakesAWhile(i int) (int, error) {
	// what this does isn’t that important for the example
	// but it’s really important in your program
	time.Sleep(2 * time.Second)
	return i * 2, nil
}

func doAnotherThing(i int) (int, error) {
	// we can wait a bit here, too
	time.Sleep(50 * time.Millisecond)
	return i + 1, nil
}

func TestCancel1(t *testing.T) {
	a := 10
	f := New(func() (interface{}, error) {
		return doSomethingThatTakesAWhile(a)
	}).Then(func(v interface{}) (interface{}, error) {
		return doAnotherThing(v.(int))
	})

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("cancelling f!")
		f.Cancel()
	}()
	val, err := f.Get()
	if val != nil || err != nil || !f.IsCancelled() {
		t.Errorf("Expected nil, nil, true, got %v, %v, %v", val, err, f.IsCancelled())
	}
}

func TestCancel2(t *testing.T) {
	a := 10
	g := New(func() (interface{}, error) {
		return doSomethingThatTakesAWhile(a)
	}).Then(func(v interface{}) (interface{}, error) {
		return doAnotherThing(v.(int))
	})

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Cancelling g (too late)!")
		g.Cancel()
	}()
	val2, err2 := g.Get()
	if val2 != 21 || err2 != nil || g.IsCancelled() {
		t.Errorf("Expected 21, nil, false, got %v, %v, %v", val2, err2, g.IsCancelled())
	}

	// once done happens, IsCancelled will never return true
	// and Get still has the calculated values
	time.Sleep(2 * time.Second)
	val2, err2 = g.Get()
	if val2 != 21 || err2 != nil || g.IsCancelled() {
		t.Errorf("Expected 21, nil, false, got %v, %v, %v", val2, err2, g.IsCancelled())
	}
}

func TestCancelConcurrent(t *testing.T) {
	loop := func() {
		const N = 8000
		start := make(chan int)
		var done sync.WaitGroup
		done.Add(N)
		f := New(func() (interface{}, error) { select {}; return 1, nil })
		var count int64
		for i := 0; i< N; i++ {
			go func() {
				defer done.Done()
				<-start
				f.Cancel()
				atomic.AddInt64(&count, 1)
			}()
		}
		close(start)
		done.Wait()
		time.Sleep(500 * time.Millisecond)
		//fmt.Println("loop done", count)
	}
	for {
		loop()
	}
}

func TestGoRoutinesQuitCancel(t *testing.T) {
	numGR := runtime.NumGoroutine()
	f := New(func() (interface{}, error) {
		time.Sleep(1 * time.Second)
		return 1, nil
	}).Then(func(i interface{}) (interface{}, error) {
		return 2, nil
	})
	f.Cancel()
	time.Sleep(2* time.Second)
	endNumGR := runtime.NumGoroutine()
	if numGR != endNumGR {
		t.Errorf("Expected the same number of goroutines at start and end, " +
			"but have %d at start and %d at end", numGR, endNumGR)
	}
}

func TestGoRoutinesQuitDone(t *testing.T) {
	numGR := runtime.NumGoroutine()
	f := New(func() (interface{}, error) {
		return 1, nil
	}).Then(func(i interface{}) (interface{}, error) {
		return 2, nil
	})
	f.Get()
	time.Sleep(100* time.Millisecond)
	endNumGR := runtime.NumGoroutine()
	if numGR != endNumGR {
		t.Errorf("Expected the same number of goroutines at start and end, " +
			"but have %d at start and %d at end", numGR, endNumGR)
	}
}