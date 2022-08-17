// Here is the core api for the hilbi shell itself
// Basically, stuff about the shell itself and other functions
// go here.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"hilbish/util"

	rt "github.com/arnodel/golua/runtime"
	"github.com/arnodel/golua/lib/packagelib"
	"github.com/maxlandon/readline"
	"github.com/blackfireio/osinfo"
	"mvdan.cc/sh/v3/interp"
)

var exports = map[string]util.LuaExport{
	"alias": {hlalias, 2, false},
	"appendPath": {hlappendPath, 1, false},
	"complete": {hlcomplete, 2, false},
	"cwd": {hlcwd, 0, false},
	"exec": {hlexec, 1, false},
	"runnerMode": {hlrunnerMode, 1, false},
	"goro": {hlgoro, 1, true},
	"highlighter": {hlhighlighter, 1, false},
	"hinter": {hlhinter, 1, false},
	"multiprompt": {hlmultiprompt, 1, false},
	"prependPath": {hlprependPath, 1, false},
	"prompt": {hlprompt, 1, true},
	"inputMode": {hlinputMode, 1, false},
	"interval": {hlinterval, 2, false},
	"read": {hlread, 1, false},
	"run": {hlrun, 1, true},
	"timeout": {hltimeout, 2, false},
	"which": {hlwhich, 1, false},
}

var hshMod *rt.Table
var hilbishLoader = packagelib.Loader{
	Load: hilbishLoad,
	Name: "hilbish",
}

func hilbishLoad(rtm *rt.Runtime) (rt.Value, func()) {
	fakeMod := rt.NewTable()
	modmt := rt.NewTable()
	mod := rt.NewTable()

	modIndex := func(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
		arg := c.Arg(1)
		val := mod.Get(arg)

		return c.PushingNext1(t.Runtime, val), nil
	}
	modNewIndex := func(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
		k, err := c.StringArg(1)
		if err != nil {
			return nil, err
		}

		v := c.Arg(2)
		if k == "highlighter" {
			var err error
			// fine to assign, since itll be either nil or a closure
			highlighter, err = c.ClosureArg(2)
			if err != nil {
				return nil, errors.New("hilbish.highlighter has to be a function")
			}
		} else if k == "hinter" {
			var err error
			hinter, err = c.ClosureArg(2)
			if err != nil {
				return nil, errors.New("hilbish.hinter has to be a function")
			}
		} else if modVal := mod.Get(rt.StringValue(k)); modVal != rt.NilValue {
			return nil, errors.New("not allowed to override in hilbish table")
		}
		mod.Set(rt.StringValue(k), v)

		return c.Next(), nil
	}
	modmt.Set(rt.StringValue("__newindex"), rt.FunctionValue(rt.NewGoFunction(modNewIndex, "__newindex", 3, false)))
	modmt.Set(rt.StringValue("__index"), rt.FunctionValue(rt.NewGoFunction(modIndex, "__index", 2, false)))
	fakeMod.SetMetatable(modmt)

	util.SetExports(rtm, mod, exports)
	hshMod = mod

	host, _ := os.Hostname()
	username := curuser.Username

	if runtime.GOOS == "windows" {
		username = strings.Split(username, "\\")[1] // for some reason Username includes the hostname on windows
	}

	util.SetFieldProtected(fakeMod, mod, "ver", rt.StringValue(getVersion()), "Hilbish version")
	util.SetFieldProtected(fakeMod, mod, "user", rt.StringValue(username), "Username of user")
	util.SetFieldProtected(fakeMod, mod, "host", rt.StringValue(host), "Host name of the machine")
	util.SetFieldProtected(fakeMod, mod, "home", rt.StringValue(curuser.HomeDir), "Home directory of the user")
	util.SetFieldProtected(fakeMod, mod, "dataDir", rt.StringValue(dataDir), "Directory for Hilbish's data files")
	util.SetFieldProtected(fakeMod, mod, "interactive", rt.BoolValue(interactive), "If this is an interactive shell")
	util.SetFieldProtected(fakeMod, mod, "login", rt.BoolValue(login), "Whether this is a login shell")
	util.SetFieldProtected(fakeMod, mod, "vimMode", rt.NilValue, "Current Vim mode of Hilbish (nil if not in Vim mode)")
	util.SetFieldProtected(fakeMod, mod, "exitCode", rt.IntValue(0), "Exit code of last exected command")
	util.Document(fakeMod, "Hilbish's core API, containing submodules and functions which relate to the shell itself.")

	// hilbish.userDir table
	hshuser := rt.NewTable()

	util.SetField(rtm, hshuser, "config", rt.StringValue(confDir), "User's config directory")
	util.SetField(rtm, hshuser, "data", rt.StringValue(userDataDir), "XDG data directory")
	util.Document(hshuser, "User directories to store configs and/or modules.")
	mod.Set(rt.StringValue("userDir"), rt.TableValue(hshuser))

	// hilbish.os table
	hshos := rt.NewTable()
	info, _ := osinfo.GetOSInfo()

	util.SetField(rtm, hshos, "family", rt.StringValue(info.Family), "Family name of the current OS")
	util.SetField(rtm, hshos, "name", rt.StringValue(info.Name), "Pretty name of the current OS")
	util.SetField(rtm, hshos, "version", rt.StringValue(info.Version), "Version of the current OS")
	util.Document(hshos, "OS info interface")
	mod.Set(rt.StringValue("os"), rt.TableValue(hshos))

	// hilbish.aliases table
	aliases = newAliases()
	aliasesModule := aliases.Loader(rtm)
	util.Document(aliasesModule, "Alias inferface for Hilbish.")
	mod.Set(rt.StringValue("aliases"), rt.TableValue(aliasesModule))

	// hilbish.history table
	historyModule := lr.Loader(rtm)
	mod.Set(rt.StringValue("history"), rt.TableValue(historyModule))
	util.Document(historyModule, "History interface for Hilbish.")

	// hilbish.completion table
	hshcomp := completionLoader(rtm)
	util.Document(hshcomp, "Completions interface for Hilbish.")
	mod.Set(rt.StringValue("completion"), rt.TableValue(hshcomp))

	// hilbish.runner table
	runnerModule := runnerModeLoader(rtm)
	util.Document(runnerModule, "Runner/exec interface for Hilbish.")
	mod.Set(rt.StringValue("runner"), rt.TableValue(runnerModule))

	// hilbish.jobs table
	jobs = newJobHandler()
	jobModule := jobs.loader(rtm)
	util.Document(jobModule, "(Background) job interface.")
	mod.Set(rt.StringValue("jobs"), rt.TableValue(jobModule))

	// hilbish.timers table
	timers = newTimerHandler()
	timerModule := timers.loader(rtm)
	util.Document(timerModule, "Timer interface, for control of all intervals and timeouts.")
	mod.Set(rt.StringValue("timers"), rt.TableValue(timerModule))

	editorModule := editorLoader(rtm)
	util.Document(editorModule, "")
	mod.Set(rt.StringValue("editor"), rt.TableValue(editorModule))

	versionModule := rt.NewTable()
	util.SetField(rtm, versionModule, "branch", rt.StringValue(gitBranch), "Git branch Hilbish was compiled from")
	util.SetField(rtm, versionModule, "full", rt.StringValue(getVersion()), "Full version info, including release name")
	util.SetField(rtm, versionModule, "commit", rt.StringValue(gitCommit), "Git commit Hilbish was compiled from")
	util.SetField(rtm, versionModule, "release", rt.StringValue(releaseName), "Release name")
	util.Document(versionModule, "Version info interface.")
	mod.Set(rt.StringValue("version"), rt.TableValue(versionModule))

	return rt.TableValue(fakeMod), nil
}

func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func setVimMode(mode string) {
	util.SetField(l, hshMod, "vimMode", rt.StringValue(mode), "Current Vim mode of Hilbish (nil if not in Vim mode)")
	hooks.Emit("hilbish.vimMode", mode)
}

func unsetVimMode() {
	util.SetField(l, hshMod, "vimMode", rt.NilValue, "Current Vim mode of Hilbish (nil if not in Vim mode)")
}

// run(cmd, returnOut) -> exitCode, stdout, stderr
// Runs `cmd` in Hilbish's sh interpreter.
// If returnOut is true, the outputs of `cmd` will be returned as the 2nd and
// 3rd values instead of being outputted to the terminal.
// --- @param cmd string
func hlrun(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	cmd, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	var terminalOut bool
	if len(c.Etc()) != 0 {
		tout := c.Etc()[0]
		termOut, ok := tout.TryBool()
		terminalOut = termOut
		if !ok {
			return nil, errors.New("bad argument to run (expected boolean, got " + tout.TypeName() + ")")
		}
	} else {
		terminalOut = true
	}

	var exitcode uint8
	stdout, stderr, err := execCommand(cmd, terminalOut)

	if code, ok := interp.IsExitStatus(err); ok {
		exitcode = code
	} else if err != nil {
		exitcode = 1
	}

	stdoutStr := ""
	stderrStr := ""
	if !terminalOut {
		stdoutStr = stdout.(*bytes.Buffer).String()
		stderrStr = stderr.(*bytes.Buffer).String()
	}

	return c.PushingNext(t.Runtime, rt.IntValue(int64(exitcode)), rt.StringValue(stdoutStr), rt.StringValue(stderrStr)), nil
}

// cwd()
// Returns the current directory of the shell
func hlcwd(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	cwd, _ := os.Getwd()

	return c.PushingNext1(t.Runtime, rt.StringValue(cwd)), nil
}


// read(prompt) -> input?
// Read input from the user, using Hilbish's line editor/input reader.
// This is a separate instance from the one Hilbish actually uses.
// Returns `input`, will be nil if ctrl + d is pressed, or an error occurs (which shouldn't happen)
// --- @param prompt string
func hlread(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	luaprompt, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	lualr := newLineReader("", true)
	lualr.SetPrompt(luaprompt)

	input, err := lualr.Read()
	if err != nil {
		return c.Next(), nil
	}

	return c.PushingNext1(t.Runtime, rt.StringValue(input)), nil
}

/*
prompt(str, typ?)
Changes the shell prompt to `str`
There are a few verbs that can be used in the prompt text.
These will be formatted and replaced with the appropriate values.
`%d` - Current working directory
`%u` - Name of current user
`%h` - Hostname of device
--- @param str string
--- @param typ string Type of prompt, being left or right. Left by default.
*/
func hlprompt(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	err := c.Check1Arg()
	if err != nil {
		return nil, err
	}
	p, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	typ := "left"
	// optional 2nd arg
	if len(c.Etc()) != 0 {
		ltyp := c.Etc()[0]
		var ok bool
		typ, ok = ltyp.TryString()
		if !ok {
			return nil, errors.New("bad argument to run (expected string, got " + ltyp.TypeName() + ")")
		}
	}

	switch typ {
		case "left":
			prompt = p
			lr.SetPrompt(fmtPrompt(prompt))
		case "right": lr.SetRightPrompt(fmtPrompt(p))
		default: return nil, errors.New("expected prompt type to be right or left, got " + typ)
	}

	return c.Next(), nil
}

// multiprompt(str)
// Changes the continued line prompt to `str`
// --- @param str string
func hlmultiprompt(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	prompt, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	multilinePrompt = prompt

	return c.Next(), nil
}

// alias(cmd, orig)
// Sets an alias of `cmd` to `orig`
// --- @param cmd string
// --- @param orig string
func hlalias(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}
	cmd, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	orig, err := c.StringArg(1)
	if err != nil {
		return nil, err
	}

	aliases.Add(cmd, orig)

	return c.Next(), nil
}

// appendPath(dir)
// Appends `dir` to $PATH
// --- @param dir string|table
func hlappendPath(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	arg := c.Arg(0)

	// check if dir is a table or a string
	if arg.Type() == rt.TableType {
		util.ForEach(arg.AsTable(), func(k rt.Value, v rt.Value) {
			if v.Type() == rt.StringType {
				appendPath(v.AsString())
			}
		})
	} else if arg.Type() == rt.StringType {
		appendPath(arg.AsString())
	} else {
		return nil, errors.New("bad argument to appendPath (expected string or table, got " + arg.TypeName() + ")")
	}

	return c.Next(), nil
}

func appendPath(dir string) {
	dir = strings.Replace(dir, "~", curuser.HomeDir, 1)
	pathenv := os.Getenv("PATH")

	// if dir isnt already in $PATH, add it
	if !strings.Contains(pathenv, dir) {
		os.Setenv("PATH", pathenv + string(os.PathListSeparator) + dir)
	}
}

// exec(cmd)
// Replaces running hilbish with `cmd`
// --- @param cmd string
func hlexec(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	cmd, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	cmdArgs, _ := splitInput(cmd)
	if runtime.GOOS != "windows" {
		cmdPath, err := exec.LookPath(cmdArgs[0])
		if err != nil {
			fmt.Println(err)
			// if we get here, cmdPath will be nothing
			// therefore nothing will run
		}

		// syscall.Exec requires an absolute path to a binary
		// path, args, string slice of environments
		syscall.Exec(cmdPath, cmdArgs, os.Environ())
	} else {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
		os.Exit(0)
	}

	return c.Next(), nil
}

// goro(fn)
// Puts `fn` in a goroutine
// --- @param fn function
func hlgoro(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	fn, err := c.ClosureArg(0)
	if err != nil {
		return nil, err
	}

	// call fn
	go func() {
		_, err := rt.Call1(l.MainThread(), rt.FunctionValue(fn), c.Etc()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error in goro function:\n\n", err)
		}
	}()

	return c.Next(), nil
}

// timeout(cb, time)
// Runs the `cb` function after `time` in milliseconds
// Returns a `timer` object (see `doc timers`).
// --- @param cb function
// --- @param time number
// --- @return table
func hltimeout(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}
	cb, err := c.ClosureArg(0)
	if err != nil {
		return nil, err
	}
	ms, err := c.IntArg(1)
	if err != nil {
		return nil, err
	}

	interval := time.Duration(ms) * time.Millisecond
	timer := timers.create(timerTimeout, interval, cb)
	timer.start()
	
	return c.PushingNext1(t.Runtime, rt.UserDataValue(timer.ud)), nil
}

// interval(cb, time)
// Runs the `cb` function every `time` milliseconds.
// Returns a `timer` object (see `doc timers`).
// --- @param cb function
// --- @param time number
// --- @return table
func hlinterval(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.CheckNArgs(2); err != nil {
		return nil, err
	}
	cb, err := c.ClosureArg(0)
	if err != nil {
		return nil, err
	}
	ms, err := c.IntArg(1)
	if err != nil {
		return nil, err
	}

	interval := time.Duration(ms) * time.Millisecond
	timer := timers.create(timerInterval, interval, cb)
	timer.start()

	return c.PushingNext1(t.Runtime, rt.UserDataValue(timer.ud)), nil
}

// complete(scope, cb)
// Registers a completion handler for `scope`.
// A `scope` is currently only expected to be `command.<cmd>`,
// replacing <cmd> with the name of the command (for example `command.git`).
// `cb` must be a function that returns a table of "completion groups."
// Check `doc completions` for more information.
// --- @param scope string
// --- @param cb function
func hlcomplete(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	scope, cb, err := util.HandleStrCallback(t, c)
	if err != nil {
		return nil, err
	}
	luaCompletions[scope] = cb

	return c.Next(), nil
}

// prependPath(dir)
// Prepends `dir` to $PATH
// --- @param dir string
func hlprependPath(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	dir, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	dir = strings.Replace(dir, "~", curuser.HomeDir, 1)
	pathenv := os.Getenv("PATH")

	// if dir isnt already in $PATH, add in
	if !strings.Contains(pathenv, dir) {
		os.Setenv("PATH", dir + string(os.PathListSeparator) + pathenv)
	}

	return c.Next(), nil
}

// which(name)
// Checks if `name` is a valid command
// --- @param binName string
func hlwhich(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	name, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	cmd := aliases.Resolve(name)

	// check for commander
	if commands[cmd] != nil {
		// they dont resolve to a path, so just send the cmd
		return c.PushingNext1(t.Runtime, rt.StringValue(cmd)), nil
	}

	path, err := exec.LookPath(cmd)
	if err != nil {
		return c.Next(), nil
	}

	return c.PushingNext1(t.Runtime, rt.StringValue(path)), nil
}

// inputMode(mode)
// Sets the input mode for Hilbish's line reader. Accepts either emacs or vim
// --- @param mode string
func hlinputMode(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	mode, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}

	switch mode {
		case "emacs":
			unsetVimMode()
			lr.rl.InputMode = readline.Emacs
		case "vim":
			setVimMode("insert")
			lr.rl.InputMode = readline.Vim
		default:
			return nil, errors.New("inputMode: expected vim or emacs, received " + mode)
	}

	return c.Next(), nil
}

// runnerMode(mode)
// Sets the execution/runner mode for interactive Hilbish. This determines whether
// Hilbish wll try to run input as Lua and/or sh or only do one of either.
// Accepted values for mode are hybrid (the default), hybridRev (sh first then Lua),
// sh, and lua. It also accepts a function, to which if it is passed one
// will call it to execute user input instead.
// --- @param mode string|function
func hlrunnerMode(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	mode := c.Arg(0)

	switch mode.Type() {
		case rt.StringType:
			switch mode.AsString() {
				case "hybrid", "hybridRev", "lua", "sh": runnerMode = mode
				default: return nil, errors.New("execMode: expected either a function or hybrid, hybridRev, lua, sh. Received " + mode.AsString())
			}
		case rt.FunctionType: runnerMode = mode
		default: return nil, errors.New("execMode: expected either a function or hybrid, hybridRev, lua, sh. Received " + mode.TypeName())
	}

	return c.Next(), nil
}

// hinter(line, pos)
// The command line hint handler. It gets called on every key insert to
// determine what text to use as an inline hint. It is passed the current
// line and cursor position. It is expected to return a string which is used
// as the text for the hint. This is by default a shim. To set hints,
// override this function with your custom handler.
// --- @param line string
// --- @param pos int
func hlhinter(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	return c.Next(), nil
}

// highlighter(line)
// Line highlighter handler. This is mainly for syntax highlighting, but in
// reality could set the input of the prompt to *display* anything. The
// callback is passed the current line and is expected to return a line that
// will be used as the input display.
// --- @param line string
func hlhighlighter(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	return c.Next(), nil
}
