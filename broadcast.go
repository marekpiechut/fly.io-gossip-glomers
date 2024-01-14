package main

type Broadcaster struct {
	store []float64
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{}
}

func (b *Broadcaster) Add(value float64) error {
	b.store = append(b.store, value)
	return nil
}
func (b *Broadcaster) Get() []float64 {
	return b.store
}
