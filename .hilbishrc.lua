-- Default Hilbish config
package.path = package.path .. ';./libs/?/init.lua' .. ';/home/sammy/.luarocks/lib/lua/5.4/?.so'

fs = require 'fs'
commander = require 'commander'

commander.register("cd", function (path)
	fs.cd(path[1])
end)
--[[commander = {
	__commands = {}
}
commander.__commands.ayo = function ()
	print("ayo?")
end]]--

local ansikit = require 'ansikit'

prompt(ansikit.text('λ {bold}{cyan}'..os.getenv('USER')..' >{magenta}>{cyan}>{reset} '))

--hook("tab complete", function ())
