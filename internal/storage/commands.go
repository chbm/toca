package commands 

type Ops int8
type Value string

const (
	Get Ops = iota
	Put Ops
	Delete Ops
)

type Command struct {
	op Ops
	key string
	value string
	r chan
}

