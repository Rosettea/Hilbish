-- Default Hilbish config
package.path = package.path .. ';./libs/?/init.lua'

local ansikit = require 'ansikit'

prompt(ansikit.text('λ {bold}{cyan}'..os.getenv('USER')..' >{magenta}>{cyan}>{reset} '))
