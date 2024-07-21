local bait = require 'bait'
local lunacolors = require 'lunacolors'

PREAMBLE = [[
Getting Started: https://rosettea.github.io/Hilbish/docs/getting-started/
🛈 These tips can be disabled with hilbish.opts.tips = false
]]

hilbish.tips = {
	"Join the discord and say hi! -> https://discord.gg/3PDdcQz",
	"{green}hilbish.alias{reset} -> Sets an alias to another cmd",
	"{green}hilbish.appendPath{reset} -> Appends the provided dir to the command path ($PATH)",
	"{green}hilbish.completions{reset} -> Are use to control suggestions when tab completing.",
	"{green}hilbish.message{reset} -> Simple notification system which can be used by other plugins and parts of the shell to notify the user of various actions.",
	[[
	{green}hilbish.opts{reset} -> Simple toggle or value options a user can set.
		- EX: hilbish.opts.greeting = false, will cause the greeting message on start-up to not display.
	]],
	[[
	{green}hilbish.runner{reset} -> The runner interface contains functions that allow the user to change how Hilbish interprets interactive input.
		- The default runners can run shell script and Lua code.
	]],
	[[
	Add Lua-written commands with the commander module!
	Checkout the docs here -> https://rosettea.github.io/Hilbish/docs/api/commander/
	]]
}

bait.catch('hilbish.init', function()
	if hilbish.interactive and hilbish.opts.tip then
		local idx = math.random(1, #hilbish.tips)
		print(lunacolors.format(PREAMBLE .. "\nTip: " .. hilbish.tips[idx]))
	end
end)
