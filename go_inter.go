package main

import "fmt"

type Shape interface {
	Sides() int
	Area() int
}

type Square struct {
	len int
}

func (s *Square) Sides() int {
	return 4
}

func (s *Square) Area() int {
	return 6
}

func main() {
	s := Square{
		len: 5,
	}
	fmt.Printf("%d\n", s.Sides())

	if( 1 != 2) {
		fmt.Print("")
	}
}

//var _ Shape = &Square{}
var _ Shape = (*Square)(nil)