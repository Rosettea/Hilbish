local commander = require 'commander'

local function removeQuote(str)
	if (string.sub(str, 1, 1) == '"' and string.sub(str, #str, #str) == '"') or
		(string.sub(str, 1, 1) == '\'' and string.sub(str, #str, #str) == '\'') then
		return string.sub(str, 2, #str - 1)
	end

	return str
end

commander.register('alias', function(args, sinks)
	local aliases = hilbish.aliases.list()
	local function printAlias(name)
		sinks.out:writeln(name .. "='" .. aliases[name] .. "'")
	end

	if #args == 0 then
		for alias, _ in pairs(aliases) do
			printAlias(alias)
		end
		return
	end

	local sepIdx = string.find(args[1], "=")
	if sepIdx == nil then
		if aliases[args[1]] == nil then
			return 1
		end

		printAlias(args[1])
		return
	end

	local alias = string.sub(args[1], 1, sepIdx - 1)
	local cmd = removeQuote(string.sub(args[1], sepIdx + 1))
	hilbish.aliases.add(alias, cmd)
end)
