local bait = require 'bait'
local lunacolors = require 'lunacolors'

PREAMBLE = [[
Getting Started: https://rosettea.github.io/Hilbish/docs/getting-started/
Documentation: https://rosettea.github.io/Hilbish/
Github: https://github.com/Rosettea/Hilbish
]]

hilbish.tips = {
	"Join the discord and say hi! -> https://discord.gg/3PDdcQz",
	"{green}hilbish.alias{reset} -> Sets an alias to another cmd",
	"{green}hilbish.appendPath{reset} -> Appends the provided dir to the command path ($PATH)",
	"{green}hilbish.appendPath{reset} -> Appends the provided dir to the command path ($PATH)",
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
