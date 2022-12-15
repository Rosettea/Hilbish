local fs = require 'fs'
local emmyPattern = '^%-%-%- (.+)'
local pieces = {}

local files = fs.readdir 'nature'
for _, fname in ipairs(files) do
	local isScript = fname:match'%.lua$'
	if not isScript then goto continue end

	local f = io.open(string.format('nature/%s', fname))
	local header = f:read '*l'
	if not header:match(emmyPattern) then goto continue end

	local iface = header:match(emmyPattern)
	pieces[iface] = {}

	local docPiece = {}

	for line in f:lines() do
		if line == header then goto continue2 end
		if not line:match(emmyPattern) then
			if line:match '^function' then
				local pattern = (string.format('^function %s.', iface) .. '(%w+)')
				local funcName = line:match(pattern)
				pieces[iface][funcName] = docPiece
			end
			docPiece = {}
			goto continue2
		end

		table.insert(docPiece, line)
		::continue2::
	end
	::continue::
end

for iface, dps in pairs(pieces) do
	local mod = iface:match '(%w+)%.'
	local path = string.format('docs/api/%s/%s.md', mod, iface)
	local f <close> = io.open(path, 'a+')

	for func, docs in pairs(dps) do
		local params = table.filter(docs, function(t)
			return t:match '^%-%-%- @param'
		end)
		f:write(string.format('## %s(', func))
		for i, str in ipairs(params) do
			if i ~= 1 then
				f:write ', '
			end
			f:write(str:match '^%-%-%- @param ([%w]+) ')
		end
		f:write(')\n')

		for _, str in ipairs(docs) do
			if not str:match '^%-%-%- @' then
				f:write(str:match '^%-%-%- (.+)' .. '\n')
			end	
		end
		f:write('\n')
	end
	f:flush()
end




