package main

import (
	"errors"
	"testing"
	"time"
)

func doSomethingThatTakesAWhile(i int) (int, error) {
	// what this does isn't that important for the example
	// but it's really important in your program
	time.Sleep(2 * time.Second)
	return i * 2, nil
}

func doAnotherThing(i interface{}) (interface{}, error) {
	ii := i.(int)
	// we can wait a bit here, too
	time.Sleep(50 * time.Millisecond)
	return ii + 1, nil
}

func TestThen(t *testing.T) {
	a := 10
	f := New(func() (interface{}, error) {
		return doSomethingThatTakesAWhile(a)
	}).Then(doAnotherThing)

	val, err := f.Get()
	if val != 21 || err != nil {
		t.Errorf("Expected 21, nil, got %v, %v", val, err)
	}
}

func doAnError(i int) (int, error) {
	//this errors out
	return 0, errors.New("Nope, no good")
}

func TestThenError(t *testing.T) {
	a := 10
	f := New(func() (interface{}, error) {
		return doAnError(a)
	}).Then(func(v interface{}) (interface{}, error) {
		return doSomethingThatTakesAWhile(v.(int))
	})

	val, err := f.Get()
	if val != 0 ||  err == nil || err.Error() != "Nope, no good" {
		t.Errorf("Expected 0, false, Nope, no good, got %v, %v", val, err)
	}
}
