local ansikit = require 'ansikit'
local bait = require 'bait'
local commander = require 'commander'
local hilbish = require 'hilbish'
local terminal = require 'terminal'
local Greenhouse = require 'nature.greenhouse'
local Page = require 'nature.greenhouse.page'

commander.register('greenhouse', function(args, sinks)
	local gh = Greenhouse(sinks.out)

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

	gh:initUi()
end)
