-- We're basically porting Ansikit to lua
-- https://github.com/Luvella/AnsiKit/blob/master/lib/index.js
-- which is made by yours truly sammy :^)
local lunacolors = require 'lunacolors'
local ansikit = {}

ansikit.clear = function(scrollback)
	typ = (scrollback and 3 or 2)
	return ansikit.printCSI(typ, 'J')
end

ansikit.clearFromPos = function(scrollback)
	return ansikit.printCSI(0, 'J')
end

ansikit.clearLine = function()
	return ansikit.printCSI(2, 'K')
end

ansikit.clearToPos = function()
	return ansikit.printCSI(1, 'J')
end

ansikit.color256 = function(color)
	color = (color and color or 0)
	return ansikit.printCSI('38;5;' .. color)
end

ansikit.cursorDown = function(y)
	y = (y and y or 1)
	return ansikit.printCSI(y, 'B')
end

ansikit.cursorLeft = function(x)
	x = (x and x or 1)
	return ansikit.printCSI(x, 'D')
end

-- TODO: cursorPos
-- https://github.com/Luvella/AnsiKit/blob/master/lib/index.js#L90

ansikit.cursorRight = function(x)
	x = (x and x or 1)
	return ansikit.printCSI(x, 'C')
end

ansikit.cursorStyle = function(style)
	style = (style and style or ansikit.underlineCursor)
	if style > 6 or style < 1 then style = ansikit.underlineCursor end
	
	return ansikit.printCSI(style, ' q')
end

ansikit.cursorTo = function(x, y)
	x, y = (x and x or 1), (y and y or 1)
	return ansikit.printCSI(x .. ';' .. y, 'H')
end

ansikit.cursorUp = function(y)
	y = (y and y or 1)
	return ansikit.printCSI(y, 'A')
end

ansikit.getCode = function(code, terminate)
	return string.char(0x001b) .. code ..
	(terminate and string.char(0x001b) .. '\\' or '')
end

ansikit.getCSI = function(code, endc)
	endc = (endc and endc or 'm')
	code = (code and code or '')

	return string.char(0x001b) .. '[' .. code .. endc
end

ansikit.hideCursor = function()
	return ansikit.printCSI('?25', 'l')
end

ansikit.link = function(url, text)
	if not url then error 'ansikit: missing url for hyperlink' end
	local text = (text and text or 'link')
	return lunacolors.blue('\27]8;;' .. url .. '\27\\' .. text .. '\27]8;;\27\\\n')
end

ansikit.print = function(text)
	io.write(ansikit.format(text))
	return ansikit
end

ansikit.printCode = function(code, terminate)
	io.write(ansikit.getCode(code, terminate))
	return ansikit
end

ansikit.printCSI = function(code, endc)
	io.write(ansikit.getCSI(code, endc))
	return ansikit
end

ansikit.println = function(text)
	print(ansikit.print(text))
	return ansikit
end

ansikit.reset = function()
	return ansikit.printCode('c')
end

ansikit.restoreCursor = function()
	return ansikit.printCSI(nil, 'u')
end

ansikit.restoreState = function()
	return ansikit.printCode(8)
end

ansikit.rgb = function(r, g, b)
	r = (r and r or 0)
	g = (g and g or 0)
	b = (b and b or 0)

	return ansikit.printCSI('38;2;' .. r .. ';' .. g .. ';' .. b)
end

ansikit.saveCursor = function()
	return ansikit.printCSI(nil, 's')
end

ansikit.saveState = function()
	return ansikit.printCode(7)
end

ansikit.setTitle = function(text)
	return ansikit.printCode(']2;' .. text, true)
end

ansikit.showCursor = function()
	return ansikit.printCSI('?25', 'h')
end

return ansikit

