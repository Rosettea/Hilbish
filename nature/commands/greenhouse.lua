local ansikit = require 'ansikit'
local bait = require 'bait'
local commander = require 'commander'
local hilbish = require 'hilbish'
local terminal = require 'terminal'
local Greenhouse = require 'nature.greenhouse'

commander.register('greenhouse', function(args, sinks)
	local fname = args[1]
	local done = false
	local f <close> = io.open(fname, 'r')
	if not f then
		sinks.err:writeln(string.format('could not open file %s', fname))
	end

	bait.catch('signal.sigint', function()
		done = true
	end)

	local gh = Greenhouse(sinks.out)
	gh:setText(f:read '*a')

	ansikit.screenAlt()
	ansikit.clear(true)
	gh:draw()

	hilbish.goro(function()
		while not done do
			local c = read()
			if c == 3 then
				done = true
			end

			if c == 27 then
				local c1 = read()
				if c1 == 91 then
					local c2 = read()
					if c2 == 66 then -- arrow down
						gh:scroll 'down'
					elseif c2 == 65 then -- arrow up
						gh:scroll 'up'
					end
				end
				goto continue
			end
			print('\nchar:')
			print(c)

			::continue::
		end
	end)

	while not done do
		--
	end
	ansikit.clear()
	ansikit.screenMain()
end)

function read()
	terminal.saveState()
	terminal.setRaw()
	local c = io.read(1)

	terminal.restoreState()
	return c:byte()
end
