package models

type Counter struct {
	Id    int
	Name  string
	Value int
}

func CreateCounter(id, value int, name string) Counter {
	return Counter{
		Id:    id,
		Value: value,
		Name:  name,
	}
}

func (c *Counter) Add(n int) {
	c.Value += n
}

func (c *Counter) Rest(n int) {
	c.Value -= n
}
