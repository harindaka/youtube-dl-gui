package main

//Counter plugin
type Counter struct {
	Value uint `json:"value"`
}

//NewCounter creates a new Counter plugin
func NewCounter() Counter {
	return Counter{
		Value: 0,
	}
}

//Add increments value
func (c *Counter) Add(value uint) {
	c.Value += value
}
