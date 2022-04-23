Hilbish has a pretty good completion system. It has a nice looking menu,
with 2 types of menus: grid (like file completions) or list.

Like most parts of Hilbish, it's made to be extensible and customizable.
The default handler for completions in general can be overwritten to provide
more advanced completions if needed.

# Completion Handler
By default, it provides 3 things: for the first argument, binaries (with a
plain name requested to complete, those in $PATH), files, or command
completions. With the default completion handler, it will try to run a
handler for the command or fallback to file completions.

To overwrite it, just assign a function to `hilbish.completion.handler`
like so:
function hilbish.completion.handler(line, pos)
	-- do things
end
It is passed 2 arguments, the entire line, and the current cursor position.
The functions in the completion interface take 3 arguments: query, ctx,
and fields. The `query`, which what the user is currently trying to complete,
`ctx`, being just the entire line, and `fields` being a table of arguments.
It's just `ctx` split up, delimited by spaces.
It's expected to return 2 things: a table of completion groups, and a prefix.
A completion group is defined as a table with 2 keys: `items` and `type`.
The `items` field is just a table of items to use for completions.
The `type` is for the completion menu type, being either `grid` or `list`.
The prefix is what all the completions start with. It should be empty
if the user doesn't have a query. If the beginning of the completion
item does not match the prefix, it will be replaced and fixed properly
in the line. It is case sensitive.

If you want to overwrite the functionality of the general completion handler,
or make your command completion have files as well (and filter them),
then there is the `files` function, which is mentioned below.

# Completion Interface
## Functions
- `files(query, ctx, fields)` -> table, prefix: get file completions, based
on the user's query.
- `bins(query, ctx, fields)` -> table, prefix: get binary/executable
completions, based on user query.
- `call(scope, query, ctx, fields)` -> table, prefix: call a completion handler
with `scope`, usually being in the form of `command.<name>`
