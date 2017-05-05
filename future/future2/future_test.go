package main

import (
	"testing"
	"time"
)

func doSomethingThatTakesAWhile(i int) (int, error) {
	// what this does isn't that important for the example
	// but it's really important in your program
	time.Sleep(2 * time.Second)
	return i * 2, nil
}

func timeFutureGet(f Future) (time.Duration, interface{}, error) {
	start := time.Now()
	val, err := f.Get()
	elapsed := time.Now().Sub(start)
	return elapsed, val, err
}

func TestGet(t *testing.T) {
	a := 10
	f := New(func() (interface{}, error) {
		return doSomethingThatTakesAWhile(a)
	})

	// this will take about 2 seconds to complete
	elapsed, val, err := timeFutureGet(f)
	if elapsed < time.Second {
		t.Errorf("That should have taken 2 seconds to finish, took %v",
			elapsed)
	}
	if val != 20 || err != nil {
		t.Errorf("Expected 20 and nil, got %v and %v", val, err)
	}

	// this will complete immediately
	elapsed, val, err = timeFutureGet(f)
	if elapsed > time.Second {
		t.Errorf("That should have taken no time to finish, took %v",
			elapsed)
	}
	if val != 20 || err != nil {
		t.Errorf("Expected 20 and nil, got %v and %v", val, err)
	}
}
