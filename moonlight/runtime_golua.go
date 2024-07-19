package moonlight

import (
	"os"

	rt "github.com/arnodel/golua/runtime"
	"github.com/arnodel/golua/lib"
	"github.com/arnodel/golua/lib/debuglib"
)

type Runtime struct{
	rt *rt.Runtime
}

func NewRuntime() *Runtime {
	r := rt.New(os.Stdout)
	r.PushContext(rt.RuntimeContextDef{
		MessageHandler: debuglib.Traceback,
	})
	lib.LoadAll(r)

	return specificRuntimeToGeneric(r)
}

func specificRuntimeToGeneric(rtm *rt.Runtime) *Runtime {
	rr := Runtime{
		rt: rtm,
	}

	return &rr
}

func (mlr *Runtime) UnderlyingRuntime() *rt.Runtime {
	return mlr.rt
}
