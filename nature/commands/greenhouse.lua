local ansikit = require 'ansikit'
local bait = require 'bait'
local commander = require 'commander'
local hilbish = require 'hilbish'
local terminal = require 'terminal'
local Greenhouse = require 'nature.greenhouse'
local Page = require 'nature.greenhouse.page'

commander.register('greenhouse', function(args, sinks)
	local gh = Greenhouse(sinks.out)
	local done = false

	if sinks['in'].pipe then
		local page = Page(sinks['in']:readAll())
		gh:addPage(page)
	end

	for _, name in ipairs(args) do
		local f <close> = io.open(name, 'r')
		if not f then
			sinks.err:writeln(string.format('could not open file %s', name))
		end

		local page = Page(f:read '*a')
		gh:addPage(page)
	end

	bait.catch('signal.sigint', function()
		done = true
	end)

	bait.catch('signal.resize', function()
		gh:update()
	end)

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

					if c2 == 49 then
						local c3 = read()
						if c3 == 59 then
							local c4 = read()
							if c4 == 53 then
								local c5 = read()
								if c5 == 67 then
									gh:next()
								elseif c5 == 68 then
									gh:previous()
								end
							end
						end
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
