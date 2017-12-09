package main

import (
	"fmt"
)

//Native plugin
type Native struct{}

//NewNative creates a new Counter plugin
func newNative() Native {
	return Native{}
}

//Add increments value
func (c *Native) Add(val1 uint, val2 uint) uint {
	fmt.Println(val1 + val2)
	return val1 + val2
}
