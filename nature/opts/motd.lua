local bait = require 'bait'
local lunacolors = require 'lunacolors'

hilbish.motd = [[
Finally at {red}v2.2!{reset} So much {green}documentation improvements{reset}
and 1 single fix for Windows! {blue}.. and a feature they can't use.{reset}
]]

bait.catch('hilbish.init', function()
	if hilbish.interactive and hilbish.opts.motd then
		print(lunacolors.format(hilbish.motd))
	end
end)
