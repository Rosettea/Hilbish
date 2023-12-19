local fs = require 'fs'
local emmyPattern = '^%-%-%- (.+)'
local modpattern = '^%-+ @module (%w+)'
local pieces = {}

local files = fs.readdir 'nature'
for _, fname in ipairs(files) do
	local isScript = fname:match'%.lua$'
	if not isScript then goto continue end

	local f = io.open(string.format('nature/%s', fname))
	local header = f:read '*l'
	local mod = header:match(modpattern)
	if not mod then goto continue end

	print(fname, mod)
	pieces[mod] = {}

	local docPiece = {}
	for line in f:lines() do
		if line == header then goto continue2 end
		if not line:match(emmyPattern) then
			if line:match '^function' then
				local pattern = (string.format('^function %s%%.', mod) .. '(%w+)')
				local funcName = line:match(pattern)
				if not funcName then goto continue2 end
				print(line)
				print(pattern)
				print(funcName)
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
	local mod = iface:match '(%w+)%.' or 'nature'
	local path = string.format('luadocs/api/%s/%s.md', mod, iface)
	local f <close> = io.open(path, 'a+')
	print(f)

	print(mod, path)
	fs.mkdir(fs.dir(path), true)

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
