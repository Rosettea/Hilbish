// Here is the core api for the hilbi shell itself
// Basically, stuff about the shell itself and other functions
// go here.
package main

import (
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
//	"github.com/blackfireio/osinfo"
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
	"prompt": {hlprompt, 1, false},
	"inputMode": {hlinputMode, 1, false},
	"interval": {hlinterval, 2, false},
	"read": {hlread, 1, false},
	"run": {hlrun, 1, false},
	"timeout": {hltimeout, 2, false},
	"which": {hlwhich, 1, false},
}

var greeting string
var hshMod *rt.Table
var hilbishLoader = packagelib.Loader{
	Load: hilbishLoad,
	Name: "hilbish",
}

func hilbishLoad(rtm *rt.Runtime) (rt.Value, func()) {
	mod := rt.NewTable()
	util.SetExports(rtm, mod, exports)
	hshMod = mod

	host, _ := os.Hostname()
	username := curuser.Username

	if runtime.GOOS == "windows" {
		username = strings.Split(username, "\\")[1] // for some reason Username includes the hostname on windows
	}

	greeting = `Welcome to {magenta}Hilbish{reset}, {cyan}` + username + `{reset}.
The nice lil shell for {blue}Lua{reset} fanatics!
Check out the {blue}{bold}guide{reset} command to get started.
`
	util.SetField(rtm, mod, "ver", rt.StringValue(version), "Hilbish version")
	util.SetField(rtm, mod, "user", rt.StringValue(username), "Username of user")
	util.SetField(rtm, mod, "host", rt.StringValue(host), "Host name of the machine")
	util.SetField(rtm, mod, "home", rt.StringValue(curuser.HomeDir), "Home directory of the user")
	util.SetField(rtm, mod, "dataDir", rt.StringValue(dataDir), "Directory for Hilbish's data files")
	util.SetField(rtm, mod, "interactive", rt.BoolValue(interactive), "If this is an interactive shell")
	util.SetField(rtm, mod, "login", rt.BoolValue(login), "Whether this is a login shell")
	util.SetField(rtm, mod, "greeting", rt.StringValue(greeting), "Hilbish's welcome message for interactive shells. It has Lunacolors formatting.")
	util.SetField(rtm, mod, "vimMode", rt.NilValue, "Current Vim mode of Hilbish (nil if not in Vim mode)")
	util.SetField(rtm, hshMod, "exitCode", rt.IntValue(0), "Exit code of last exected command")
	util.Document(mod, "Hilbish's core API, containing submodules and functions which relate to the shell itself.")

	// hilbish.userDir table
	hshuser := rt.NewTable()

	util.SetField(rtm, hshuser, "config", rt.StringValue(confDir), "User's config directory")
	util.SetField(rtm, hshuser, "data", rt.StringValue(userDataDir), "XDG data directory")
	util.Document(hshuser, "User directories to store configs and/or modules.")
	mod.Set(rt.StringValue("userDir"), rt.TableValue(hshuser))

/*
	// hilbish.os table
	hshos := L.NewTable()
	info, _ := osinfo.GetOSInfo()

	util.SetField(L, hshos, "family", lua.LString(info.Family), "Family name of the current OS")
	util.SetField(L, hshos, "name", lua.LString(info.Name), "Pretty name of the current OS")
	util.SetField(L, hshos, "version", lua.LString(info.Version), "Version of the current OS")
	util.Document(L, hshos, "OS info interface")
	L.SetField(mod, "os", hshos)
*/

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
	hshcomp := rt.NewTable()
	util.SetField(rtm, hshcomp, "files",
	rt.FunctionValue(rt.NewGoFunction(luaFileComplete, "files", 3, false)),
	"Completer for files")

	util.SetField(rtm, hshcomp, "bins",
	rt.FunctionValue(rt.NewGoFunction(luaBinaryComplete, "bins", 3, false)),
	"Completer for executables/binaries")

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

	return rt.TableValue(mod), nil
}

func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func luaFileComplete(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	query, ctx, fds, err := getCompleteParams(t, c)
	if err != nil {
		return nil, err
	}

	completions := fileComplete(query, ctx, fds)
	luaComps := rt.NewTable()

	for i, comp := range completions {
		luaComps.Set(rt.IntValue(int64(i + 1)), rt.StringValue(comp))
	}

	return c.PushingNext1(t.Runtime, rt.TableValue(luaComps)), nil
}

func luaBinaryComplete(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	query, ctx, fds, err := getCompleteParams(t, c)
	if err != nil {
		return nil, err
	}

	completions, _ := binaryComplete(query, ctx, fds)
	luaComps := rt.NewTable()

	for i, comp := range completions {
		luaComps.Set(rt.IntValue(int64(i + 1)), rt.StringValue(comp))
	}

	return c.PushingNext1(t.Runtime, rt.TableValue(luaComps)), nil
}

func getCompleteParams(t *rt.Thread, c *rt.GoCont) (string, string, []string, error) {
	if err := c.CheckNArgs(3); err != nil {
		return "", "", []string{}, err
	}
	query, err := c.StringArg(0)
	if err != nil {
		return "", "", []string{}, err
	}
	ctx, err := c.StringArg(1)
	if err != nil {
		return "", "", []string{}, err
	}
	fields, err := c.TableArg(2)
	if err != nil {
		return "", "", []string{}, err
	}

	var fds []string
	nextVal := rt.NilValue
	for {
		next, val, ok := fields.Next(nextVal)
		if next == rt.NilValue {
			break
		}
		nextVal = next

		valStr, ok := val.TryString()
		if !ok {
			continue
		}

		fds = append(fds, valStr)
	}

	return query, ctx, fds, err
}

func setVimMode(mode string) {
	util.SetField(l, hshMod, "vimMode", rt.StringValue(mode), "Current Vim mode of Hilbish (nil if not in Vim mode)")
	hooks.Em.Emit("hilbish.vimMode", mode)
}

func unsetVimMode() {
	util.SetField(l, hshMod, "vimMode", rt.NilValue, "Current Vim mode of Hilbish (nil if not in Vim mode)")
}

// run(cmd)
// Runs `cmd` in Hilbish's sh interpreter.
// --- @param cmd string
func hlrun(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	cmd, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	var exitcode uint8
	err = execCommand(cmd)

	if code, ok := interp.IsExitStatus(err); ok {
		exitcode = code
	} else if err != nil {
		exitcode = 1
	}

	return c.PushingNext1(t.Runtime, rt.IntValue(int64(exitcode))), nil
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
prompt(str)
Changes the shell prompt to `str`
There are a few verbs that can be used in the prompt text.
These will be formatted and replaced with the appropriate values.
`%d` - Current working directory
`%u` - Name of current user
`%h` - Hostname of device
--- @param str string
*/
func hlprompt(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	var prompt string
	err := c.Check1Arg()
	if err == nil {
		prompt, err = c.StringArg(0)
	}
	if err != nil {
		return nil, err
	}
	lr.SetPrompt(fmtPrompt(prompt))

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
		nextVal := rt.NilValue
		for {
			next, val, ok := arg.AsTable().Next(nextVal)
			if next == rt.NilValue {
				break
			}
			nextVal = next

			valStr, ok := val.TryString()
			if !ok {
				continue
			}

			appendPath(valStr)
		}
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
// --- @param cb function
// --- @param time number
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

	timeout := time.Duration(ms) * time.Millisecond
	time.Sleep(timeout)

	_, err = rt.Call1(l.MainThread(), rt.FunctionValue(cb)) 
	if err != nil {
		return nil, err
	}

	return c.Next(), nil
}

// interval(cb, time)
// Runs the `cb` function every `time` milliseconds
// --- @param cb function
// --- @param time number
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

	ticker := time.NewTicker(interval)
	stop := make(chan rt.Value)

	go func() {
		for {
			select {
			case <-ticker.C:
				_, err := rt.Call1(l.MainThread(), rt.FunctionValue(cb)) 
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error in interval function:\n\n", err)
					stop <- rt.BoolValue(true) // stop the interval
				}
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	// TODO: return channel
	return c.Next(), nil
}

// complete(scope, cb)
// Registers a completion handler for `scope`.
// A `scope` is currently only expected to be `command.<cmd>`,
// replacing <cmd> with the name of the command (for example `command.git`).
// `cb` must be a function that returns a table of "completion groups."
// A completion group is a table with the keys `items` and `type`.
// `items` being a table of items and `type` being the display type of
// `grid` (the normal file completion display) or `list` (with a description)
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

// which(binName)
// Searches for an executable called `binName` in the directories of $PATH
// --- @param binName string
func hlwhich(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	binName, err := c.StringArg(0)
	if err != nil {
		return nil, err
	}
	path, err := exec.LookPath(binName)
	if err != nil {
		return c.Next(), nil
	}

	return c.PushingNext1(t.Runtime, rt.StringValue(path)), nil
}

// inputMode(mode)
// Sets the input mode for Hilbish's line reader. Accepts either emacs for vim
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
				// no fallthrough doesnt work so eh
				case "hybrid", "hybridRev", "lua", "sh": runnerMode = mode
				default: return nil, errors.New("execMode: expected either a function or hybrid, hybridRev, lua, sh. Received " + mode.AsString())
			}
		case rt.FunctionType: runnerMode = mode
		default: return nil, errors.New("execMode: expected either a function or hybrid, hybridRev, lua, sh. Received " + mode.TypeName())
	}

	return c.Next(), nil
}

// hinter(cb)
// Sets the hinter function. This will be called on every key insert to determine
// what text to use as an inline hint. The callback is passed 2 arguments:
// the current line and the position. It is expected to return a string
// which will be used for the hint.
// --- @param cb function
func hlhinter(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	hinterCb, err := c.ClosureArg(0)
	if err != nil {
		return nil, err
	}
	hinter = hinterCb
	
	return c.Next(), err
}

// highlighter(cb)
// Sets the highlighter function. This is mainly for syntax hightlighting, but in
// reality could set the input of the prompt to display anything. The callback
// is passed the current line as typed and is expected to return a line that will
// be used to display in the line.
// --- @param cb function
func hlhighlighter(t *rt.Thread, c *rt.GoCont) (rt.Cont, error) {
	if err := c.Check1Arg(); err != nil {
		return nil, err
	}
	highlighterCb, err := c.ClosureArg(0)
	if err != nil {
		return nil, err
	}
	highlighter = highlighterCb

	return c.Next(), err
}
