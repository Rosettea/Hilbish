package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"

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
	once bool
	args []string
	// save path for a few reasons, one being security (lmao) while the other
	// would just be so itll be the same binary command always (path changes)
	path string
	handle *exec.Cmd
	cmdout io.Writer
	cmderr io.Writer
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func (j *job) start() error {
	if j.handle == nil || j.once {
		// cmd cant be reused so make a new one
		cmd := exec.Cmd{
			Path: j.path,
			Args: j.args,
		}
		j.setHandle(&cmd)
	}
	// bgProcAttr is defined in execfile_<os>.go, it holds a procattr struct
	// in a simple explanation, it makes signals from hilbish (sigint)
	// not go to it (child process)
	j.handle.SysProcAttr = bgProcAttr
	// reset output buffers
	j.stdout.Reset()
	j.stderr.Reset()
	// make cmd write to both standard output and output buffers for lua access
	j.handle.Stdout = io.MultiWriter(j.cmdout, j.stdout)
	j.handle.Stderr = io.MultiWriter(j.cmderr, j.stderr)

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

func (j *job) wait() {
	j.handle.Wait()
}

func (j *job) setHandle(handle *exec.Cmd) {
	j.handle = handle
	j.args = handle.Args
	j.path = handle.Path
	if handle.Stdout != nil {
		j.cmdout = handle.Stdout
	}
	if handle.Stderr != nil {
		j.cmderr = handle.Stderr
	}
}

func (j *job) getProc() *os.Process {
	handle := j.handle
	if handle != nil {
		return handle.Process
	}

	return nil
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
	luaJob.Set(rt.StringValue("stdout"), rt.StringValue(string(j.stdout.Bytes())))
	luaJob.Set(rt.StringValue("stderr"), rt.StringValue(string(j.stderr.Bytes())))

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
		cmdout: os.Stdout,
		cmderr: os.Stderr,
		stdout: &bytes.Buffer{},
		stderr: &bytes.Buffer{},
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

func (j *jobHandler) disown(id int) error {
	j.mu.RLock()
	if j.jobs[id] == nil {
		return errors.New("job doesnt exist")
	}
	j.mu.RUnlock()

	j.mu.Lock()
	delete(j.jobs, id)
	j.mu.Unlock()

	return nil
}

func (j *jobHandler) stopAll() {
	j.mu.RLock()
	defer j.mu.RUnlock()

	for _, jb := range j.jobs {
		// on exit, unix shell should send sighup to all jobs
		if jb.running {
			proc := jb.getProc()
			proc.Signal(syscall.SIGHUP)
			jb.wait() // waits for program to exit due to sighup
		}
	}
}

func (j *jobHandler) loader(rtm *rt.Runtime) *rt.Table {
	jobFuncs := map[string]util.LuaExport{
		"all": {j.luaAllJobs, 0, false},
		"last": {j.luaLastJob, 0, false},
		"get": {j.luaGetJob, 1, false},
		"add": {j.luaAddJob, 2, false},
		"disown": {j.luaDisownJob, 1, false},
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

func (j *jobHandler) luaDisownJob(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	jobID, err := c.IntArg(0)
	if err != nil {
		return nil, err
	}

	err = j.disown(int(jobID))
	if err != nil {
		return nil, err
	}

	return c.Next(), nil
}

func (j *jobHandler) luaLastJob(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	job := j.jobs[j.latestID]
	if job == nil { // incase we dont have any jobs yet
		return c.Next(), nil
	}

	return c.PushingNext1(t.Runtime, job.lua()), nil
}
