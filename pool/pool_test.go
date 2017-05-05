package pool

import (
	"sync"
	"testing"
	"fmt"
	"time"
)

func adderFactory() interface{} {
	return func(a, b int) int {
		time.Sleep(1 * time.Second)
		return a + b
	}
}

func TestPool(t *testing.T) {
	adderPool := NewPool(adderFactory, 3)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				//borrow an adder
				adder := adderPool.Borrow().(func(int, int) int)
				fmt.Println("adding on goroutine", i)
				fmt.Println("result == ", adder(i, i + 1))
				//return the item
				adderPool.Return(adder)
			}
		}(i)
	}
	wg.Wait()
}
