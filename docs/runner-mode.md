Hilbish is *unique,* when interactive it first attempts to run input as
Lua and then tries shell script. But if you're normal, you wouldn't
really be using Hilbish anyway but you'd also not want this
(or maybe want Lua only in some cases.)

The "runner mode" of Hilbish is customizable via `hilbish.runnerMode`,
which determines how Hilbish will run user input. By default, this is
set to `hybrid` which is the previously mentioned behaviour of running Lua
first then going to shell script. If you want the reverse order, you can
set it to `hybridRev` and for isolated modes there is `sh` and `lua`
respectively.

You can also set it to a function, which will be called everytime Hilbish
needs to run interactive input. For example, you can set this to a simple
function to compile and evaluate Fennel, and now you can run Fennel.
You can even mix it with sh to make a hybrid mode with Lua replaced by
Fennel.

An example:
hilbish.runnerMode(function(input)
	local ok = pcall(fennel.eval, input)
	if ok then
		return input, 0, nil
	end

	return hilbish.runner.sh(input)
end)

The `hilbish.runner` interface is an alternative to using `hilbish.runnerMode`
and also provides the sh and Lua runner functions that Hilbish itself uses.
A runner function is expected to return 3 values: the input, exit code, and an error.
The input return is there incase you need to prompt for more input.
If you don't, just return the input passed to the runner function.
The exit code has to be a number, it will be 0 otherwise and the error can be
`nil` to indicate no error.

## Functions
These are the "low level" functions for the `hilbish.runner` interface.

+ setMode(mode) > The same as `hilbish.runnerMode`
+ sh(input) -> table > Runs `input` in Hilbish's sh interpreter
+ lua(input) -> table > Evals `input` as Lua code

The table value that runners return can have at least 4 values:
+ input (string): The full input text.
+ exitCode (number): Exit code (usually from a command)
+ continue (boolean): Whether to prompt the user for more input
(in the case of incomplete syntax)
+ err (string): A string that represents an error from the runner.
This should only be set when, for example, there is a syntax error.
It can be set to a few special values for Hilbish to throw the right
hooks and have a better looking message.

+ `<command>: not-found` will throw a `command.not-found` hook
based on what `<command>` is.
+ `<command>: not-executable` will throw a `command.not-executable` hook.

The others here are defined in Lua and have EmmyLua documentation.
These functions should be preferred over the previous ones.
+ setCurrent(mode) > The same as `setMode`, but works with runners managed
via the functions below.
+ add(name, runner) > Adds a runner to a table of available runners. The `runner`
argument is either a function or a table with a run callback.
+ set(name, runner) > The same as `add` but requires passing a table and
overwrites if the `name`d runner already exists.
+ get(name) > runner > Gets a runner by name. It is a table with at least a
run function, to run input.
+ exec(cmd, runnerName) > Runs `cmd` with a runner. If `runnerName` isn't passed,
the current runner mode is used.
