local bait = require 'bait'
local lunacolors = require 'lunacolors'

hilbish.motd = [[
1000 commits on the Hilbish repository brings us to {cyan}Version 2.1!{reset}
Docs, docs, docs... At least builtins work with pipes now.
]]

bait.catch('hilbish.init', function()
	if hilbish.interactive and hilbish.opts.motd then
		print(lunacolors.format(hilbish.motd))
	end
end)
