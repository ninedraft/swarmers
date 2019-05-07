package swarm

import (
	"math"

	 "github.com/atedja/go-vector"
)

func (drone *Drone) move(until float64) {
	if until <= 0 {
		until = 1
	}
	var coords = vector.NewWithValues(drone.CoordsV())
	var targetPoint =drone.targetPoint
	var dd = vector.NewWithValues([]float64{
		drone.Velocity() * drone.rnd.Float64(),
		drone.Velocity() * drone.rnd.Float64(),
	})
	dd.DoWithIndex(func(i int, v float64) float64 {
		return v*sign(coords[i] - targetPoint[i])
	})
	var newCoords = vector.Add(coords, dd)
	drone.body.Move(newCoords[0], newCoords[1])
	if vector.Subtract(newCoords, targetPoint).Magnitude() < until {
		drone.state = Patrol
	}
}

func sign(x float64) float64 {
	if math.Signbit(x) {
		return -1
	}
	return 1
}
