package main

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
	"github.com/arnodel/golua/lib/iolib"
)

var jobs *jobHandler

type job struct {
	cmd string
	running bool
	id int
	pid int
	exitCode int
	once bool
	args []string
	// save path for a few reasons, one being security (lmao) while the other
	// would just be so itll be the same binary command always (path changes)
	path string
	handle *exec.Cmd
	stdin io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (j *job) start() error {
	if j.handle == nil || j.once {
		// cmd cant be reused so make a new one
		cmd := exec.Cmd{
			Path: j.path,
			Args: j.args,
			Stdin: j.stdin,
			Stdout: j.stdout,
			Stderr: j.stderr,
		}
		j.setHandle(&cmd)
	}

	if !j.once {
		j.once = true
	}

	err := j.handle.Start()
	proc := j.getProc()

	j.pid = proc.Pid
	j.running = true

	hooks.Em.Emit("job.start", j.lua())

	return err
}

func (j *job) stop() {
	// finish will be called in exec handle
	proc := j.getProc()
	if proc != nil {
		proc.Kill()
	}
}

func (j *job) finish() {
	j.running = false
	hooks.Em.Emit("job.done", j.lua())
}

func (j *job) setHandle(handle *exec.Cmd) {
	j.handle = handle
	j.args = handle.Args
	j.path = handle.Path
	j.stdin = handle.Stdin
	j.stdout = handle.Stdout
	j.stderr = handle.Stderr
}

func (j *job) getProc() *os.Process {
	handle := j.handle
	if handle != nil {
		return handle.Process
	}

	return nil
}

func (j *job) setStdio(typ string, f *iolib.File) {
	switch typ {
		case "in": j.stdin = f.File
		case "out": j.stdout = f.File
		case "err": j.stderr = f.File
	}
}

func (j *job) lua() rt.Value {
	jobFuncs := map[string]util.LuaExport{
		"stop": {j.luaStop, 0, false},
		"start": {j.luaStart, 0, false},
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

func (j *job) luaStart(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if !j.running {
		err := j.start()
		exit := handleExecErr(err)
		j.exitCode = int(exit)
		j.finish()
	}

	return c.Next(), nil
}

func (j *job) luaStop(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if j.running {
		j.stop()
		j.finish()
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

func (j *jobHandler) add(cmd string, args []string, path string) *job {
	j.mu.Lock()
	defer j.mu.Unlock()

	j.latestID++
	jb := &job{
		cmd: cmd,
		running: false,
		id: j.latestID,
		args: args,
		path: path,
		stdin: os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	j.jobs[j.latestID] = jb
	hooks.Em.Emit("job.add", jb.lua())

	return jb
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
		"add": {j.luaAddJob, 2, false},
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

func (j *jobHandler) luaAddJob(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}
	cmd, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	largs, err := c.TableArg(1)
	if err != nil {
		return nil, err
	}

	var args []string
	util.ForEach(largs, func(k rt.Value, v rt.Value) {
		if v.Type() == rt.StringType {
			args = append(args, v.AsString())
		}
	})
	// TODO: change to lookpath for args[0]
	jb := j.add(cmd, args, args[0])

	return c.PushingNext1(t.Runtime, jb.lua()), nil
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
