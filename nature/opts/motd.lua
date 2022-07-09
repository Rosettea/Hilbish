local bait = require 'bait'
local lunacolors = require 'lunacolors'

hilbish.motd = [[
Hilbish 2.0 is a {red}major{reset} update! If your config doesn't work
anymore, that will definitely be why! A MOTD, very message, much day.
]]

bait.catch('hilbish.init', function()
	if hilbish.opts.motd then
		print(lunacolors.format(hilbish.motd))
	end
end)
