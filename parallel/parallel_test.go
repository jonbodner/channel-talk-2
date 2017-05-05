package parallel

import (
	"fmt"
	"testing"
)

type in struct {
	a int
	b int
}

type out struct {
	source string
	result int
}

func adder(i in) (out, error) {
	r := i.a + i.b
	return out{"adder", r}, nil
}

func multiplier(i in) (out, error) {
	r := i.a * i.b
	return out{"multiplier", r}, nil
}

func subtracter(i in) (out, error) {
	r := i.a - i.b
	return out{"subtracter", r}, nil
}

func divider(i in) (out, error) {
	r := i.a / i.b
	return out{"divider", r}, nil
}

var evaluators = []Evaluator{
	func(inV interface{}) (interface{}, error) {
		return adder(inV.(in))
	},
	func(inV interface{}) (interface{}, error) {
		return multiplier(inV.(in))
	},
	func(inV interface{}) (interface{}, error) {
		return subtracter(inV.(in))
	},
	func(inV interface{}) (interface{}, error) {
		return divider(inV.(in))
	},
}

func TestFanout(t *testing.T) {
	results, errors := FanOut(evaluators, in{2, 3})
	fmt.Println("results:", results)
	if len(results) != 4 {
		t.Errorf("Expected 4 results, got %d", len(results))
	}
	fmt.Println("errors:", errors)
	if len(errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(errors))
	}
}
