package main

//Counter plugin
type Counter struct {
	value uint
}

//NewCounter creates a new Counter plugin
func NewCounter() Counter {
	return Counter{
		value: 0,
	}
}

func (c *Counter) add(value uint) {
	c.value += value
}
