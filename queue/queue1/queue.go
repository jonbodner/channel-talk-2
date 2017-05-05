package main

func MakeInfiniteQueue() Queue {
	in := make(chan interface{})
	out := make(chan interface{})
	//do magic stuff 
	return &queueInner{in: in, out: out}
}
