package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"hilbish/util"
	
	rt "github.com/arnodel/golua/runtime"
)

type timerType int64
const (
	timerInterval timerType = iota
	timerTimeout
)

type timer struct{
	id int
	typ timerType
	running bool
	dur time.Duration
	fun *rt.Closure
	th *timerHandler
	ticker *time.Ticker
	channel chan struct{}
}

func (t *timer) start() error {
	if t.running {
		return errors.New("timer is already running")
	}

	t.running = true
	t.th.running++
	t.ticker = time.NewTicker(t.dur)

	go func() {
		for {
			select {
			case <-t.ticker.C:
				_, err := rt.Call1(l.MainThread(), rt.FunctionValue(t.fun))
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error in function:\n", err)
					t.stop()
				}
				// only run one for timeout
				if t.typ == timerTimeout {
					t.stop()
				}
			case <-t.channel:
				t.ticker.Stop()
				return
			}
		}
	}()
	
	return nil
}

func (t *timer) stop() error {
	if !t.running {
		return errors.New("timer not running")
	}

	t.channel <- struct{}{}
	t.running = false
	t.th.running--
	
	return nil
}

func (t *timer) luaStart(thr *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	err := t.start()
	if err != nil {
		return nil, err
	}
	
	return c.Next(), nil
}

func (t *timer) luaStop(thr *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	err := t.stop()
	if err != nil {
		return nil, err
	}
	
	return c.Next(), nil
}

func (t *timer) lua() rt.Value {
	tExports := map[string]util.LuaExport{
		"start": {t.luaStart, 0, false},
		"stop": {t.luaStop, 0, false},
	}
	luaTimer := rt.NewTable()
	util.SetExports(l, luaTimer, tExports)

	luaTimer.Set(rt.StringValue("type"), rt.IntValue(int64(t.typ)))
	luaTimer.Set(rt.StringValue("running"), rt.BoolValue(t.running))
	luaTimer.Set(rt.StringValue("duration"), rt.IntValue(int64(t.dur / time.Millisecond)))

	return rt.TableValue(luaTimer)
}
