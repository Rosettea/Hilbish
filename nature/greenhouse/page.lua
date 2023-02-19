local Object = require 'nature.object'

local Page = Object:extend()

function Page:new(text)
	self:setText(text)
end

function Page:setText(text)
	self.lines = string.split(text, '\n')
end

return Page
