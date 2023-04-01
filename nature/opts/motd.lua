local bait = require 'bait'
local lunacolors = require 'lunacolors'

hilbish.motd = [[
1000 commits on the Hilbish repository brings us to {cyan}Version 2.1!{reset}
Docs, docs, docs... At least builtins work with pipes now.
]]

bait.catch('hilbish.init', function()
	if hilbish.interactive and hilbish.opts.motd then
		if os.date '%d' == '01' and os.date '%m' == '04' then
			print('lolololololololol\n')
		else
			print(lunacolors.format(hilbish.motd))
		end
	end
end)
