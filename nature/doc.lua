local lunacolors = require 'lunacolors'

local M = {}

function M.highlight(text)
	return text:gsub('\'.-\'', lunacolors.yellow)
	--:gsub('%-%- .-', lunacolors.black)
end

function M.renderCodeBlock(text)
	local longest = 0
	local lines = string.split(text:gsub('\t', '    '), '\n')

	for i, line in ipairs(lines) do
		local len = line:len()
		if len > longest then longest = len end
	end

	for i, line in ipairs(lines) do
		lines[i] = M.highlight(line:sub(0, longest))
		.. string.rep(' ', longest - line:len())
	end

	return '\n' .. lunacolors.format('{greyBg}' .. table.concat(lines, '\n')) .. '\n'
end

return M
