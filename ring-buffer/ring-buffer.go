package ringbuffer

import (
	"fmt"
	"log"
)

// A channel-based ring buffer removes the oldest item when the queue is full
// Ref:
// https://tanzu.vmware.com/content/blog/a-channel-based-ring-buffer-in-go

func NewRingBuffer(inCh, outCh chan int) *ringBuffer {
	return &ringBuffer{
		inCh:  inCh,
		outCh: outCh,
	}
}

// ringBuffer throttle buffer for implement async channel.
type ringBuffer struct {
	inCh  chan int
	outCh chan int
}

func (r *ringBuffer) Run() {

	for v := range r.inCh {
		select {
		case r.outCh <- v:
		default:
			fmt.Println("default running")
			<-r.outCh // pop one item from outchan
			r.outCh <- v
		}
	}
	close(r.outCh)
}

func Run() {
	inCh := make(chan int)
	outCh := make(chan int, 10) // try to change outCh buffer to understand the result
	rb := NewRingBuffer(inCh, outCh)
	go rb.Run()

	for i := 0; i < 10; i++ {
		inCh <- i
	}

	close(inCh)

	for res := range outCh {
		log.Println(res)
	}

}
