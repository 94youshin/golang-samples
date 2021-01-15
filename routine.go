package main

import (
	"fmt"
	"time"
)

var eventCh chan int

func main() {
	eventCh = make(chan int)
	go productor(eventCh)
	mid()

	time.Sleep(15 * time.Second)
}

func mid() {
	go resumer(eventCh)
	return
}

func productor(eventCh chan<- int ) {
	for i:=0; i<10; i++ {
		eventCh <- i
		time.Sleep(1 * time.Second)
	}
}

func resumer(eventCh <- chan int) {
	for i := range eventCh {
		fmt.Println(i)
		if i == 9 {
			return
		}
	}
}