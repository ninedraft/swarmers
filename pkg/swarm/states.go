package swarm

type State uint64

const (
	Patrol State = 1 + iota
	Moving
)
