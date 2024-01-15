package main

type Broadcaster struct {
	store []int
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{}
}

func (b *Broadcaster) Add(value int) error {
	b.store = append(b.store, value)
	return nil
}
func (b *Broadcaster) Get() []int {
	return b.store
}
