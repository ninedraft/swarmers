package swarm

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/ninedraft/r"
)

type Port interface {
	Sender
	ID() uint64
	Alive() bool
}

type void struct{}

func (void) Send(Message) {}
func (void) ID() uint64   { return 0 }
func (void) Alive() bool  { return true }

type Sender interface {
	Send(msg Message)
}

var _ Port = &Drone{}
var _ Sender = &Drone{}

type Drone struct {
	alive  uint64
	id     uint64
	queue  chan Message
	ports  []Port
	rnd    *rand.Rand
	memory [256]float64
	state  State

	targetPoint []float64
	body        *body
}

type DroneConfig struct {
	ID              uint64
	X, Y            float64
	Friends         []*Drone
	RndSeed         int64
	QueueBufferSize uint64
	Stored          float64
	Energy          uint64
}

func NewDrone(config DroneConfig) *Drone {
	var ports = make([]Port, 0, len(config.Friends))
	for _, friend := range config.Friends {
		ports = append(ports, friend)
	}

	var seed = config.RndSeed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	var rnd = rand.New(rand.NewSource(seed))

	var QueueBufferSize = config.QueueBufferSize
	if QueueBufferSize == 0 {
		QueueBufferSize = 256
	}
	return &Drone{
		id:    config.ID,
		queue: make(chan Message, QueueBufferSize),
		ports: ports,
		rnd:   rnd,

		body: &body{
			x:  config.X,
			y:  config.Y,
			st: config.Stored,
			e:  config.Energy,
		},
	}
}

// threadsafe
func (drone *Drone) ID() uint64 {
	return drone.id
}

// threadsafe
func (drone *Drone) Alive() bool {
	return atomic.LoadUint64(&drone.alive) == 0
}

// threadsafe
func (drone *Drone) markDead() {
	atomic.StoreUint64(&drone.alive, 666)
}

// threadsafe
func (drone *Drone) self() SenderID {
	return SenderID(drone.id)
}

func (drone *Drone) Run(after ...func()) {
	defer func() {
		for _, hook := range after {
			hook()
		}
	}()
	for drone.Alive() {
		select {
		case msg, notClosed := <-drone.queue:
			if _, ok := msg.(Die); ok || !notClosed {
				drone.markDead()
				return
			}
			drone.handleMessage(msg)
		default:
			drone.routine()
		}
	}
}

// threadsafe
func (drone *Drone) Send(msg Message) {
	if msg != nil && drone.Alive() {
		drone.queue <- msg
	}
}

func (drone *Drone) friendsN() int {
	return len(drone.ports)
}

func (drone *Drone) choicePort() Port {
	var n = drone.friendsN()
	var inds = r.R(n).Ints()
	drone.rnd.Shuffle(n, func(i, j int) {
		inds[i], inds[j] = inds[j], inds[i]
	})
	for _, i := range inds {
		if drone.ports[i].Alive() {
			return drone.ports[i]
		}
	}
	return void{}
}

func (drone *Drone) routine() {
	if drone.body.e == 0 {
		drone.markDead()
	}
	switch drone.state {
	case Moving:
		drone.move(0)
	default:
		drone.state = Patrol
	}
}

func (drone *Drone) handleMessage(msg Message) {
	switch msg := msg.(type) {
	case Remember:
		drone.memory[msg.Key] = msg.Value
	case Retrieve:
		drone.choicePort().Send(Remember{
			SenderID: drone.self(),
			Value:    drone.memory[msg.Key],
		})
	case Sleep:
		time.Sleep(msg.Duration)
	case MoveTo:
		drone.targetPoint = []float64{msg.X, msg.Y}
		drone.state = Moving
	default:
		// pass
	}
}

// threadsafe
func (drone *Drone) CoordsV() []float64 {
	return drone.body.CoordsV()
}

// threadsafe
func (drone *Drone) Coords() (x, y float64) {
	return drone.body.Coords()
}

// threadsafe
func (drone *Drone) Velocity() float64 {
	return drone.body.Velocity()
}

// threadsafe
func (drone *Drone) String() string {
	return fmt.Sprintf("Drone %d: %s", drone.id, drone.body)
}
