package main

import (
	"log"
	"math/rand"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/maps"
)

type Broadcaster struct {
	mutex    sync.Mutex
	node     *maelstrom.Node
	store    map[int]bool
	topology topology
}

func NewBroadcaster(node *maelstrom.Node) *Broadcaster {
	return &Broadcaster{
		node:  node,
		store: make(map[int]bool),
		topology: topology{
			node: node,
		},
	}
}

func (b *Broadcaster) Add(value int) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if !b.store[value] {
		b.store[value] = true
		if err := b.propagate(value); err != nil {
			log.Printf("Broadcast: failed to propagate after %d", value)
		}
	}
	return nil
}

func (b *Broadcaster) Get() []int {
	return maps.Keys(b.store)
}

type PropagateBody struct {
	Value int `json:"value"`
}

func (b *Broadcaster) propagate(value int) error {
	friends := b.topology.friends
	for _, friend := range friends {
		b.node.Send(friend, map[string]any{
			"type":  "propagate",
			"value": value,
		})
	}
	return nil
}

type topology struct {
	node    *maelstrom.Node
	friends []string
}

func (t *topology) Update(neighbourhood map[string][]string) {
	var friends []string
	if passedNodes, ok := neighbourhood[t.node.ID()]; ok {
		friends = append(friends, passedNodes...)
	}

	if len(friends) == 0 {
		log.Printf("Broadcast: Empty friends, pick a random other node")
		allNodes := t.node.NodeIDs()
		if len(allNodes) > 1 {
			randomNode := rand.Intn(len(allNodes) - 1)
			for i, e := range allNodes {
				if i >= randomNode {
					friends = []string{e}
					break
				}
			}
		} else {
			log.Printf("Broadcast: No other nodes found, nowhere to broadcast")
		}
	}
	log.Printf("Broadcast: Update topology on %s - %v", t.node.ID(), friends)
	t.friends = friends
}
