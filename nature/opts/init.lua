local opts = {}
hilbish.opts = {}

setmetatable(hilbish.opts, {
	__newindex = function(_, k, v)
		if opts[k] == nil then
			error(string.format('opt %s does not exist', k))
		end

		opts[k] = v
	end,
	__index = function(_, k)
		return opts[k]
	end
})

local function setupOpt(name, default)
	opts[name] = default
	require('nature.opts.' .. name)
end

local defaultOpts = {
	autocd = false,
	greeting = string.format([[Welcome to {magenta}Hilbish{reset}, {cyan}%s{reset}.
The nice lil shell for {blue}Lua{reset} fanatics!
]], hilbish.user)
}

for optsName, default in pairs(defaultOpts) do
	setupOpt(optsName, default)
end
