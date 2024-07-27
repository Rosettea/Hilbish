local bait = require 'bait'
local lunacolors = require 'lunacolors'

hilbish.motd = [[
Wait ... {magenta}2.3{reset} is basically the same as {red}2.2?{reset}
Erm.. {blue}Ctrl-C works for Commanders,{reset} {cyan}and the sh runner has some fixes.{reset}
Just trust me bro, this is an imporant bug fix release. {red}- 🌺 sammyette{reset}
]]

bait.catch('hilbish.init', function()
	if hilbish.interactive and hilbish.opts.motd then
		print(lunacolors.format(hilbish.motd))
	end
end)
