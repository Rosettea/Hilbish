//go:build midnight

package moonlight

type Callable interface {
	Continuation(*Runtime, Cont) Cont
}

type Closure struct {
	refIdx int // so since we cant store the actual lua closure,
	// we need a index to the ref in the lua registry... or something like that.
}

func (mlr *Runtime) ClosureArg(num int) (*Closure, error) {
	return &Closure{
		refIdx: -1,
	}, nil
}

/*
func (c *Closure) Continuation(mlr *Runtime, c Cont) Cont {
}
*/
