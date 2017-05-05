package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func doSomethingThatTakesAWhile(i int) (int, error) {
	// what this does isn’t that important for the example
	// but it’s really important in your program
	time.Sleep(2 * time.Second)
	fmt.Println("finished doSomethingThatTakesAWhile")
	return i * 2, nil
}

func doAnotherThing(i interface{}) (interface{}, error) {
	ii := i.(int)
	// we can wait a bit here, too
	time.Sleep(50 * time.Millisecond)
	fmt.Println("finished doAnotherThing")
	return ii + 1, nil
}

func TestCancel1(t *testing.T) {
	a := 10
	f := New(func() (interface{}, error) {
		return doSomethingThatTakesAWhile(a)
	}).Then(doAnotherThing)

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
	f := New(func() (interface{}, error) {
		return doSomethingThatTakesAWhile(a)
	}).Then(doAnotherThing)

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Cancelling g (too late)!")
		f.Cancel()
	}()
	val, err := f.Get()
	if val != 21 || err != nil || f.IsCancelled() {
		t.Errorf("Expected 21, nil, false, got %v, %v, %v", val, err, f.IsCancelled())
	}
}

func TestCancelConcurrent(t *testing.T) {
	const N = 8000
	loop := func() {
		start := make(chan int)
		var done sync.WaitGroup
		done.Add(N)
		f := New(func() (interface{}, error) { select {}; return 1, nil })
		for i := 0; i< N; i++ {
			go func() {
				defer done.Done()
				<-start
				f.Cancel()
			}()
		}
		close(start)
		done.Wait()
	}
	for {
		loop()
	}
}
