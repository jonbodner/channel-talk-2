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

func timeFutureGetUntil(f Future, d time.Duration) (time.Duration, interface{}, error) {
	start := time.Now()
	val, err := f.GetUntil(d)
	elapsed := time.Now().Sub(start)
	return elapsed, val, err
}

func TestGetUntil(t *testing.T) {
	a := 10
	f := New(func() (interface{}, error) {
		return doSomethingThatTakesAWhile(a)
	})

	elapsed, val, err := timeFutureGetUntil(f, time.Second)
	if elapsed < time.Second {
		t.Errorf("That should have waited a second to finish, took %v", elapsed)
	}
	if val != nil || err != FUTURE_TIMEOUT {
		t.Errorf("Expected nil, and FUTURE_TIMEOUT, got %v and %v", val, err)
	}

	elapsed, val, err = timeFutureGet(f)
	if elapsed > 1100 * time.Millisecond {
		t.Errorf("That should have taken about a second to finish, took %v",
			elapsed)
	}
	if val != 20 || err != nil {
		t.Errorf("Expected 20 and nil, got %v and %v", val, err)
	}

	elapsed, val, err = timeFutureGetUntil(f, time.Second)
	if elapsed > 100 * time.Millisecond {
		t.Errorf("That should have taken no time at all to finish, took %v",
			elapsed)
	}
	if val != 20 || err != nil {
		t.Errorf("Expected 20 and nil, got %v and %v", val, err)
	}
}
