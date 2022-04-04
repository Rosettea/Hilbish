package main

import (
	"sync"
	"os"

	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
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

func (j *job) lua() rt.Value {
	jobFuncs := map[string]util.LuaExport{
		"stop": {j.luaStop, 0, false},
	}
	luaJob := rt.NewTable()
	util.SetExports(l, luaJob, jobFuncs)

	luaJob.Set(rt.StringValue("cmd"), rt.StringValue(j.cmd))
	luaJob.Set(rt.StringValue("running"), rt.BoolValue(j.running))
	luaJob.Set(rt.StringValue("id"), rt.IntValue(int64(j.id)))
	luaJob.Set(rt.StringValue("pid"), rt.IntValue(int64(j.pid)))
	luaJob.Set(rt.StringValue("exitCode"), rt.IntValue(int64(j.exitCode)))

	return rt.TableValue(luaJob)
}

func (j *job) luaStop(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if j.running {
		j.stop()
	}

	return c.Next(), nil
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

func (j *jobHandler) loader(rtm *rt.Runtime) *rt.Table {
	jobFuncs := map[string]util.LuaExport{
		"all": {j.luaAllJobs, 0, false},
		"get": {j.luaGetJob, 1, false},
	}

	luaJob := rt.NewTable()
	util.SetExports(rtm, luaJob, jobFuncs)

	return luaJob
}

func (j *jobHandler) luaGetJob(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	jobID, err := c.IntArg(0)
	if err != nil {
		return nil, err
	}

	job := j.jobs[int(jobID)]
	if job == nil {
		return c.Next(), nil
	}

	return c.PushingNext1(t.Runtime, job.lua()), nil
}

func (j *jobHandler) luaAllJobs(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	jobTbl := rt.NewTable()
	for id, job := range j.jobs {
		jobTbl.Set(rt.IntValue(int64(id)), job.lua())
	}

	return c.PushingNext1(t.Runtime, rt.TableValue(jobTbl)), nil
}
