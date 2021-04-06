-- We're basically porting Ansikit to lua
-- https://github.com/Luvella/AnsiKit/blob/master/lib/index.js
-- which is made by yours truly sammy :^)

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

ansikit.format = function(text)
	local colors = {
		-- TODO: write codes manually instead of using functions
		-- less function calls = faster ????????
		reset = {'{reset}', ansikit.getCSI(0)},
		bold = {'{bold}', ansikit.getCSI(1)},
		dim = {'{dim}', ansikit.getCSI(2)},
		italic = {'{italic}', ansikit.getCSI(3)},
		underline = {'{underline}', ansikit.getCSI(4)},
		invert = {'{invert}', ansikit.getCSI(7)},
		bold_off = {'{bold-off}', ansikit.getCSI(22)},
		underline_off = {'{underline-off}', ansikit.getCSI(24)},
		black = {'{black}', ansikit.getCSI(30)},
		red = {'{red}', ansikit.getCSI(31)},
		green = {'{green}', ansikit.getCSI(32)},
		yellow = {'{yellow}', ansikit.getCSI(33)},
		blue = {'{blue}', ansikit.getCSI(34)},
		magenta = {'{magenta}', ansikit.getCSI(35)},
		cyan = {'{cyan}', ansikit.getCSI(36)},
		white = {'{white}', ansikit.getCSI(37)},
		red_bg = {'{red-bg}', ansikit.getCSI(41)},
		green_bg = {'{green-bg}', ansikit.getCSI(42)},
		yellow_bg = {'{green-bg}', ansikit.getCSI(43)},
		blue_bg = {'{blue-bg}', ansikit.getCSI(44)},
		magenta_bg = {'{magenta-bg}', ansikit.getCSI(45)},
		cyan_bg = {'{cyan-bg}', ansikit.getCSI(46)},
		white_bg = {'{white-bg}', ansikit.getCSI(47)},
		gray = {'{gray}', ansikit.getCSI(90)},
		bright_red = {'{bright-red}', ansikit.getCSI(91)},
		bright_green = {'{bright-green}', ansikit.getCSI(92)},
		bright_yellow = {'{bright-yellow}', ansikit.getCSI(93)},
		bright_blue = {'{bright-blue}', ansikit.getCSI(94)},
		bright_magenta = {'{bright-magenta}', ansikit.getCSI(95)},
		bright_cyan = {'{bright-cyan}', ansikit.getCSI(96)}
	}

	for k, v in pairs(colors) do
		text = text:gsub(v[1], v[2])
	end

	return text
end

ansikit.getCode = function(code, terminate)
	return string.char 0x001b .. code ..
	(terminate and string.char 0x001b .. '\\' or '')
end

ansikit.getCSI = function(code, endc)
	endc = (endc and endc or 'm')
	code = (code and code or '')

	return string.char 0x001b .. '[' .. code .. endc
end

ansikit.hideCursor = function()
	return ansikit.printCSI('?25', 'l')
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
	return ansikit.printCode 'c'
end

ansikit.restoreCursor = function()
	return ansikit.printCSI(nil, 'u')
end

ansikit.restoreState = function()
	return ansikit.printCode 8
end

ansikit.rgb = function(r, g, b)
	r = (r and r or 0)
	g = (g and g or 0)
	b = (b and b or 0)

	return ansikit.printCSI '38;2;' .. r .. ';' .. g .. ';' .. b
end

ansikit.saveCursor = function()
	return ansikit.printCSI(nil, 's')
end

ansikit.saveState = function()
	return ansikit.printCode 7
end

ansikit.setTitle = function(text)
	ansikit.printCode (']2;' .. text, true)
end

ansikit.showCursor = function()
	return ansikit.printCSI('?25', 'h')
end

return ansikit

