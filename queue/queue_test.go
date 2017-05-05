package queue_test

import (
	"testing"
	"time"
	"sync"
	"github.com/jonbodner/channel-talk-2/queue"
)

func TestWriteValuesNoPause(t *testing.T) {
	queue := queue.MakeInfiniteQueue()
	writeVals(queue, 0)
}

func TestMakeInfiniteNoPause(t *testing.T) {
	doTest(0, 0, t)
}

func TestMakeInfiniteSlowWrite(t *testing.T) {
	doTest(0, 50*time.Millisecond, t)
}

func TestMakeInfiniteSlowRead(t *testing.T) {
	doTest(50*time.Millisecond, 0, t)
}

func doTest(readDelay time.Duration, writeDelay time.Duration, t *testing.T) {
	queue := queue.MakeInfiniteQueue()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		writeVals(queue, writeDelay)
		wg.Done()
	}()
	go func() {
		readVals(queue, readDelay, t)
		wg.Done()
	}()
	wg.Wait()
}

func readVals(queue queue.Queue, readDelay time.Duration, t *testing.T) {
	lastVal := -1
	for v, ok := queue.Get(); ok; v, ok = queue.Get() {
		time.Sleep(readDelay)
		vi := v.(int)
		//fmt.Println(vi)
		if lastVal+1 != vi {
			t.Errorf("Unexpected value; expected %d, got %d", lastVal+1, vi)
		}
		lastVal = vi
	}
	//fmt.Println("finished reading")

	if lastVal != 99 {
		t.Errorf("Didn't get all values, last one received was %d", lastVal)
	}
}

func writeVals(queue queue.Queue, writeDelay time.Duration) {
	for i := 0; i < 100; i++ {
		//fmt.Println("writing", i)
		queue.Put(i)
		time.Sleep(writeDelay)
	}
	queue.Close()
	//fmt.Println("finished writing")
}
