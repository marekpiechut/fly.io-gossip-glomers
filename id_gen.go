package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type IdGenerator interface {
	Next() string
}

type UUIDIdGenerator struct {
}

func NewUUIDIdGenerator() *UUIDIdGenerator {

	return &UUIDIdGenerator{}
}

func (g *UUIDIdGenerator) Next() string {
	return uuid.Must(uuid.NewRandom()).String()
}

type SnowflakeIdGenerator struct {
	mutex   sync.Mutex
	seqTime int64
	seq     int64
	cluster int
	node    int
	epoch   int64
}

const SNOWFLAKE_EPOCH = 1288834974657

func NewSnowflakeIdGenerator(cluster int, node int) *SnowflakeIdGenerator {
	return NewSnowflakeIdGeneratorWithEpoch(SNOWFLAKE_EPOCH, cluster, node)
}

func NewSnowflakeIdGeneratorWithEpoch(epoch int64, cluster int, node int) *SnowflakeIdGenerator {
	if cluster < 0 {
		log.Print("Snowflake: using random cluster ID. Collisions may happen! You should provide cluster ID")
		cluster = rand.Int()
	}

	if node < 0 {
		log.Print("Snowflake: using random node ID. Collisions may happen! You should always provide node ID")
		node = rand.Int()
	}

	return &SnowflakeIdGenerator{
		cluster: cluster,
		node:    node,
		epoch:   epoch,
	}
}

func (g *SnowflakeIdGenerator) Next() string {
	var seq int64
	var now int64

	func() {
		//We need to make sure we update sequence atomically
		g.mutex.Lock()
		defer g.mutex.Unlock()

		now = time.Now().UnixMilli()
		if g.seqTime != now {
			g.seqTime = now
			g.seq = 0
		} else {
			g.seq++
		}
		seq = g.seq
	}()

	tid := ((now - g.epoch) & 0x1FFFFFFFFFF) << 22
	cid := (int64(g.cluster) & 0x1F) << 15
	nid := (int64(g.cluster) & 0x1F) << 10
	sid := (seq & 0xFFF)

	id := tid + cid + nid + sid

	return fmt.Sprintf("%d", id)
}
