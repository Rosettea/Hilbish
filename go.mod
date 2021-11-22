module hilbish

go 1.16

require (
	github.com/Rosettea/Hilbiline v0.0.0-20210603231612-80054dac3650
	github.com/Rosettea/readline v0.0.0-20211122152601-6d95ce44b7ed
	github.com/bobappleyard/readline v0.0.0-20150707195538-7e300e02d38e
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pborman/getopt v1.1.0
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9
	golang.org/x/term v0.0.0-20210916214954-140adaaadfaf
	layeh.com/gopher-luar v1.0.8
	mvdan.cc/sh/v3 v3.3.0
)

replace mvdan.cc/sh/v3 => github.com/Rosettea/sh/v3 v3.4.0-0.dev.0.20211022004519-f67a49cb50f5
