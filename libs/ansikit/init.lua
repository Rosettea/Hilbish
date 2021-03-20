-- We're basically porting Ansikit to lua
-- https://github.com/Luvella/AnsiKit/blob/master/lib/index.js
-- which is made by yours truly sammy :^)

local ansikit = {}

ansikit.getCSI = function (code, endc)
	endc = (endc and endc or 'm')
	return string.char(0x001b) .. '[' .. code .. endc
end

ansikit.text = function (text)
	local colors = {
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
		cyan = {'{cyan}', ansikit.getCSI(36)}
		-- TODO: Background, bright colors
	}

	for k, v in pairs(colors) do
		text = text:gsub(v[1], v[2])
	end

	return text
end

return ansikit

