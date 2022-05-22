package main

import (
	"sync"
	"time"

	"hilbish/util"
	
	rt "github.com/arnodel/golua/runtime"
)

var timers *timerHandler
type timerHandler struct {
	mu *sync.RWMutex
	wg *sync.WaitGroup
	timers map[int]*timer
	latestID int
	running int
}

func newTimerHandler() *timerHandler {
	return &timerHandler{
		timers: make(map[int]*timer),
		latestID: 0,
		mu: &sync.RWMutex{},
		wg: &sync.WaitGroup{},
	}
}

func (th *timerHandler) wait() {
	th.wg.Wait()
}

func (th *timerHandler) create(typ timerType, dur time.Duration, fun *rt.Closure) *timer {
	th.mu.Lock()
	defer th.mu.Unlock()

	th.latestID++
	t := &timer{
		typ: typ,
		fun: fun,
		dur: dur,
		channel: make(chan struct{}, 1),
		th: th,
		id: th.latestID,
	}
	th.timers[th.latestID] = t
	
	return t
}

func (th *timerHandler) get(id int) *timer {
	th.mu.RLock()
	defer th.mu.RUnlock()

	return th.timers[id]
}

func (th *timerHandler) luaCreate(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(3); err != nil {
		return nil, err
	}
	timerTypInt, err := c.IntArg(0)
	if err != nil {
		return nil, err
	}
	ms, err := c.IntArg(1)
	if err != nil {
		return nil, err
	}
	cb, err := c.ClosureArg(2)
	if err != nil {
		return nil, err
	}

	timerTyp := timerType(timerTypInt)
	tmr := th.create(timerTyp, time.Duration(ms) * time.Millisecond, cb)
	return c.PushingNext1(t.Runtime, tmr.lua()), nil
}

func (th *timerHandler) luaGet(thr *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	id, err := c.IntArg(0)
	if err != nil {
		return nil, err
	}

	t := th.get(int(id))
	if t != nil {
		return c.PushingNext1(thr.Runtime, t.lua()), nil
	}

	return c.Next(), nil
}

func (th *timerHandler) loader(rtm *rt.Runtime) *rt.Table {
	thExports := map[string]util.LuaExport{
		"create": {th.luaCreate, 3, false},
		"get": {th.luaGet, 1, false},
	}

	luaTh := rt.NewTable()
	util.SetExports(rtm, luaTh, thExports)

	return luaTh
}
