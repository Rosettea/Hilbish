module hilbish

go 1.17

require (
	github.com/blackfireio/osinfo v1.0.3
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9
	github.com/maxlandon/readline v0.1.0-beta.0.20211027085530-2b76cabb8036
	github.com/pborman/getopt v1.1.0
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	layeh.com/gopher-luar v1.0.10
	mvdan.cc/sh/v3 v3.4.3
)

replace mvdan.cc/sh/v3 => github.com/Rosettea/sh/v3 v3.4.0-0.dev.0.20220306140409-795a84b00b4e

replace github.com/maxlandon/readline => ./readline

replace layeh.com/gopher-luar => github.com/layeh/gopher-luar v1.0.10
