package main

import (
	"github.com/google/uuid"
)

type IdGenerator struct {
}

func NewIdGenerator() *IdGenerator {

	return &IdGenerator{}
}

func (g *IdGenerator) Next() string {
	return uuid.Must(uuid.NewRandom()).String()
}
