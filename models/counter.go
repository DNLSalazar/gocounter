package models

type Counter struct {
	Id    int64
	Name  string
	Value int
}

func CreateCounter(id int64, value int, name string) Counter {
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

type CounterCreate struct {
	Name  string
	Value int
}

func CreateCounterBody(name string, value int) CounterCreate {
	return CounterCreate{
		Name:  name,
		Value: value,
	}
}
