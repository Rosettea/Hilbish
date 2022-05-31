package main

import (
	"fmt"
	"sync"
	"time"

	"hilbish/util"
	
	rt "github.com/arnodel/golua/runtime"
)

var timers *timerHandler
var timerMetaKey = rt.StringValue("hshtimer")

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
	t.ud = timerUserData(t)

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
	return c.PushingNext1(t.Runtime, rt.UserDataValue(tmr.ud)), nil
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
		return c.PushingNext1(thr.Runtime, rt.UserDataValue(t.ud)), nil
	}

	return c.Next(), nil
}

func (th *timerHandler) loader(rtm *rt.Runtime) *rt.Table {
	timerMethods := rt.NewTable()
	timerFuncs := map[string]util.LuaExport{
		"start": {timerStart, 1, false},
		"stop": {timerStop, 1, false},
	}
	util.SetExports(rtm, timerMethods, timerFuncs)

	timerMeta := rt.NewTable()
	timerIndex := func(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
		ti, _ := timerArg(c, 0)

		arg := c.Arg(1)
		val := timerMethods.Get(arg)

		if val != rt.NilValue {
			return c.PushingNext1(t.Runtime, val), nil
		}

		keyStr, _ := arg.TryString()

		switch keyStr {
			case "type": val = rt.IntValue(int64(ti.typ))
			case "running": val = rt.BoolValue(ti.running)
			case "duration": val = rt.IntValue(int64(ti.dur / time.Millisecond))
		}

		return c.PushingNext1(t.Runtime, val), nil
	}

	timerMeta.Set(rt.StringValue("__index"), rt.FunctionValue(rt.NewGoFunction(timerIndex, "__index", 2, false)))
	l.SetRegistry(timerMetaKey, rt.TableValue(timerMeta))

	thExports := map[string]util.LuaExport{
		"create": {th.luaCreate, 3, false},
		"get": {th.luaGet, 1, false},
	}

	luaTh := rt.NewTable()
	util.SetExports(rtm, luaTh, thExports)

	return luaTh
}

func timerArg(c *rt.GoCont, arg int) (*timer, error) {
	j, ok := valueToTimer(c.Arg(arg))
	if !ok {
		return nil, fmt.Errorf("#%d must be a timer", arg + 1)
	}

	return j, nil
}

func valueToTimer(val rt.Value) (*timer, bool) {
	u, ok := val.TryUserData()
	if !ok {
		return nil, false
	}

	j, ok := u.Value().(*timer)
	return j, ok
}

func timerUserData(j *timer) *rt.UserData {
	timerMeta := l.Registry(timerMetaKey)
	return rt.NewUserData(j, timerMeta.AsTable())
}
