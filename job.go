package main

import (
	"sync"

	"github.com/yuin/gopher-lua"
)

var jobs *jobHandler

type job struct {
	cmd string
	running bool
	id int
	pid int
	exitCode int
}

func (j *job) start(pid int) {
	j.pid = pid
	j.running = true
	hooks.Em.Emit("job.start", j.lua())
}

func (j *job) finish() {
	j.running = false
	hooks.Em.Emit("job.done", j.lua())
}

func (j *job) lua() *lua.LTable {
	// returns lua table for job
	// because userdata is gross
	luaJob := l.NewTable()

	l.SetField(luaJob, "cmd", lua.LString(j.cmd))
	l.SetField(luaJob, "running", lua.LBool(j.running))
	l.SetField(luaJob, "id", lua.LNumber(j.id))
	l.SetField(luaJob, "pid", lua.LNumber(j.pid))
	l.SetField(luaJob, "exitCode", lua.LNumber(j.exitCode))

	return luaJob
}

type jobHandler struct {
	jobs map[int]*job
	latestID int
	mu *sync.RWMutex
}

func newJobHandler() *jobHandler {
	return &jobHandler{
		jobs: make(map[int]*job),
		latestID: 0,
		mu: &sync.RWMutex{},
	}
}

func (j *jobHandler) add(cmd string) {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.latestID++
	j.jobs[j.latestID] = &job{
		cmd: cmd,
		running: false,
		id: j.latestID,
	}
}

func (j *jobHandler) getLatest() *job {
	j.mu.RLock()
	defer j.mu.RUnlock()

	return j.jobs[j.latestID]
}
