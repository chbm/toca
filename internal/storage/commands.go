package storage 

type Ops int8
type Value struct {
	V string
	Exists bool
}


const (
	Get Ops = iota
	Put Ops = iota
	Delete Ops = iota
)

type Command struct {
	Op Ops
	Key string
	Value string
	R chan Value
}

