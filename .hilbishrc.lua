-- Default Hilbish config
package.path = package.path .. ';./libs/?/init.lua;/usr/share/hilbish/libs/?/init.lua'

fs = require 'fs'
commander = require 'commander'

commander.register("cd", function (path)
	if path then
		fs.cd(path[1])
	end
end)

local ansikit = require 'ansikit'

prompt(ansikit.text('λ {bold}{cyan}'..os.getenv('USER')..' >{magenta}>{cyan}>{reset} '))

--hook("tab complete", function ())
