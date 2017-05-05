package main

func MakeInfiniteQueue() Queue {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		var inQueue []interface{}
		curVal := func() interface{} {
			if len(inQueue) == 0 {
				return nil
			}
			return inQueue[0]
		}
		for in != nil {
			select {
			case v, ok := <-in:
				if !ok {
					in = nil
				} else {
					inQueue = append(inQueue, v)
				}
			case out <- curVal():
				if len(inQueue) > 0 {
					inQueue = inQueue[1:]
				}
			}
		}
		close(out)
	}()
	return &queueInner{in, out}
}
