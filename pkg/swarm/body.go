package swarm

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type body struct {
	sync.RWMutex
	x, y float64
	v    float64
	st   float64
	e    uint64
}

func (body *body) Velocity() float64 {
	return body.v
}

func (body *body) Move(coords ...float64) []float64 {
	var old = body.CoordsV()
	body.Lock()
	defer body.Unlock()
	body.x, body.y = coords[0], coords[1]
	return old
}

func (body *body) Coords() (x, y float64) {
	body.RLock()
	defer body.RUnlock()
	return body.x, body.y
}

func (body *body) CoordsV() []float64 {
	var x, y = body.Coords()
	return []float64{x, y}
}

func (body *body) Energy() uint64 {
	return atomic.LoadUint64(&body.e)
}

func (body *body) SetEnergy(energy uint64) uint64 {
	return atomic.SwapUint64(&body.e, energy)
}

func (body *body) Stored() float64 {
	body.RLock()
	defer body.RUnlock()
	return body.st
}

func (body *body) Store(value float64) float64 {
	body.Lock()
	defer body.Unlock()
	var storage = body.st
	body.st = value
	return storage
}

func (body *body) String() string {
	var x, y = body.Coords()
	var energy = body.Energy()
	var storage = body.Stored()
	return fmt.Sprintf("[%.2f, %.2f], ⚡️ %d, 📦 %.2f", x, y, energy, storage)
}
