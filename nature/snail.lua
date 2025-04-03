local hilbish = require 'hilbish'
local snail = require 'snail'

hilbish.snail = snail.new()

--- run(cmd, streams) -> exitCode (number), stdout (string), stderr (string)
--- Runs `cmd` in Hilbish's shell script interpreter.
--- The `streams` parameter specifies the output and input streams the command should use.
--- For example, to write command output to a sink.
--- As a table, the caller can directly specify the standard output, error, and input
--- streams of the command with the table keys `out`, `err`, and `input` respectively.
--- As a boolean, it specifies whether the command should use standard output or return its output streams.
--- #param cmd string
--- #param streams table|boolean
--- #returns number, string, string
--- #example
--- This code is the same as `ls -l | wc -l`
--- local fs = require 'fs'
--- local pr, pw = fs.pipe()
--- hilbish.run('ls -l', {
--- 	stdout = pw,
--- 	stderr = pw,
--- })
--- pw:close()
--- hilbish.run('wc -l', {
--- 	stdin = pr
--- })
--- #example
function hilbish.run(cmd, streams)
	local sinks = {}

	if type(streams) == 'boolean' then
		if not streams then
			sinks = {
				out = hilbish.sink.new(),
				err = hilbish.sink.new(),
				input = io.stdin
			}
		end
	elseif type(streams) == 'table' then
		sinks = streams
	end

	local out = hilbish.snail:run(cmd, {sinks = sinks})
	local returns = {out}

	if type(streams) == 'boolean' and not streams then
		table.insert(returns, sinks.out:readAll())
		table.insert(returns, sinks.err:readAll())
	end

	return table.unpack(returns)
end
