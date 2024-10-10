package main

import (
	"fmt"
)

type DummyStruct struct {
	Name  string
	Value int
}

func Dummy() {
	fmt.Println("Dummy")
}

func CreateStruct(n string, v int) *DummyStruct {
	t := DummyStruct{
		Name:  n,
		Value: v,
	}
	return &t
}

func CopyStruct(addr *DummyStruct) (res DummyStruct) {
	res = DummyStruct{
		Name:  addr.Name,
		Value: addr.Value,
	}
	return
}

func Recieve(v bool) {
	V := <-results
}

func main() {
	results := make(chan bool, 3)

}
