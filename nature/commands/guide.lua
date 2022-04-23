local ansikit = require 'ansikit'
local commander = require 'commander'

local helpTexts = {
[[
Hello there! Welcome to Hilbish, the comfy and nice little shell for
Lua users and fans. Hilbish is configured with Lua, and its
scripts are also in Lua. It also runs both Lua and shell script when
interactive (aka normal usage).
]],
[[
What does that mean for you, the user? It means that if you prefer to
use Lua for scripting instead of shell script but still have ordinary
shell usage for interactive use.
]],
[[
If this is your first time using Hilbish and Lua, check out the
Programming in Lua book here: https://www.lua.org/pil
After (or if you already know Lua) check out the doc command.
It is an in shell tool for documentation about Hilbish provided
functions and modules.
]],
[[
If you've updated from a pre-1.0 version (0.7.1 as an example)
you'll want to move your config from ~/.hilbishrc.lua to
]] ..
hilbish.userDir.config .. '/hilbish/init.lua' ..
[[

and also change all global functions (prompt, alias) to be
in the hilbish module (hilbish.prompt, hilbish.alias as examples).

And if this is your first time (most likely), you can copy a config
from ]] .. hilbish.dataDir,
[[
Since 1.0 is a big release, you'll want to check the changelog
at https://github.com/Rosettea/Hilbish/releases/tag/v1.0.0
to find more breaking changes.
]]
}
commander.register('guide', function()
	ansikit.clear()
	ansikit.cursorTo(0, 0)
	for _, text in ipairs(helpTexts) do
		print(text)
		local out = hilbish.read('Hit enter to continue ')
		ansikit.clear()
		ansikit.cursorTo(0, 0)
		if not out then
			return
		end
	end
	print 'Hope you enjoy using Hilbish!'
end)
