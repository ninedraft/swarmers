package swarm

import (
	"time"
)

type Message interface {
	From() uint64
}

type SenderID uint64

func (id SenderID) From() uint64 {
	return uint64(id)
}

type Die struct {
	SenderID
}

type Remember struct {
	SenderID
	Key   byte
	Value float64
}

type Retrieve struct {
	SenderID
	Key byte
}

type Reduce struct {
	SenderID
	Keys   []byte
	Target float64
	Op     func([]float64) float64
}

type MoveTo struct {
	SenderID
	X, Y float64
}

type Sleep struct {
	SenderID
	Duration time.Duration
}
