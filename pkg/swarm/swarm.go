package swarm

import "math/rand"

type Swarm struct {
	rnd    *rand.Rand
	drones []*Drone
	logs   chan string
}
