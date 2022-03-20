package main

import (
	"sync"
	"os"

	"github.com/yuin/gopher-lua"
)

var jobs *jobHandler

type job struct {
	cmd string
	running bool
	id int
	pid int
	exitCode int
	proc *os.Process
}

func (j *job) start(pid int) {
	j.pid = pid
	j.running = true
	hooks.Em.Emit("job.start", j.lua())
}

func (j *job) stop() {
	// finish will be called in exec handle
	j.proc.Kill()
}

func (j *job) finish() {
	j.running = false
	hooks.Em.Emit("job.done", j.lua())
}

func (j *job) setHandle(handle *os.Process) {
	j.proc = handle
}

func (j *job) lua() *lua.LTable {
	// returns lua table for job
	// because userdata is gross
	jobFuncs := map[string]lua.LGFunction{
		"stop": j.luaStop,
	}
	luaJob := l.SetFuncs(l.NewTable(), jobFuncs)

	l.SetField(luaJob, "cmd", lua.LString(j.cmd))
	l.SetField(luaJob, "running", lua.LBool(j.running))
	l.SetField(luaJob, "id", lua.LNumber(j.id))
	l.SetField(luaJob, "pid", lua.LNumber(j.pid))
	l.SetField(luaJob, "exitCode", lua.LNumber(j.exitCode))

	return luaJob
}

func (j *job) luaStop(L *lua.LState) int {
	if j.running {
		j.stop()
	}

	return 0
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


func (j *jobHandler) loader(L *lua.LState) *lua.LTable {
	jobFuncs := map[string]lua.LGFunction{
		"all": j.luaAllJobs,
		"get": j.luaGetJob,
	}

	luaJob := l.SetFuncs(l.NewTable(), jobFuncs)

	return luaJob
}

func (j *jobHandler) luaGetJob(L *lua.LState) int {
	j.mu.RLock()
	defer j.mu.RUnlock()

	jobID := L.CheckInt(1)
	job := j.jobs[jobID]
	if job != nil {
		return 0
	}
	L.Push(job.lua())

	return 1
}

func (j *jobHandler) luaAllJobs(L *lua.LState) int {
	j.mu.RLock()
	defer j.mu.RUnlock()

	jobTbl := L.NewTable()
	for id, job := range j.jobs {
		jobTbl.Insert(id, job.lua())
	}

	L.Push(jobTbl)
	return 1
}
