package storage 

type Ops int8
type Errs int8
type Value struct {
	V string
	Exists bool
}


const (
	Get Ops = iota
	Put Ops = iota
	Delete Ops = iota
	CreateNs Ops = iota
	SaveNs Ops = iota
	LoadNs Ops = iota
)

const (
	Success Errs = iota
	Failure Errs = iota
	Conflict Errs = iota
	NoNS Errs = iota
)

type Command struct {
	Op Ops
	Ns string
	Key string
	Value string
	R chan Result
}

type Result struct {
	Val Value
	Err Errs
}
