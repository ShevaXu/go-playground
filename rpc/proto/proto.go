package proto

import (
	"sync"
)

type NoArgs struct{}

type NoReturn struct{}

type Counter struct {
	mu    sync.Mutex
	count int
}

func (c *Counter) Add(delta int, ret *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count += delta
	*ret = c.count
	return nil
}

func (c *Counter) Get(arg NoArgs, ret *int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	*ret = c.count
	return nil
}

func (c *Counter) Clear(arg NoArgs, ret *NoArgs) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count = 0
	return nil
}

func NewCounter() *Counter {
	return &Counter{}
}
