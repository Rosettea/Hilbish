module hilbish

go 1.16

require (
	github.com/Rosettea/Hilbiline v0.0.0-20210624011007-8088a2d84b65
	github.com/Rosettea/readline v0.0.0-20211206201906-5570232ed3b3
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9
	github.com/creack/goselect v0.1.2 // indirect
	github.com/creack/termios v0.0.0-20160714173321-88d0029e36a1 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/maxlandon/readline v0.1.0-beta.0.20211027085530-2b76cabb8036
	github.com/pborman/ansi v1.0.0 // indirect
	github.com/pborman/getopt v1.1.0
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	layeh.com/gopher-luar v1.0.10
	mvdan.cc/sh/v3 v3.4.1
)

replace mvdan.cc/sh/v3 => github.com/Rosettea/sh/v3 v3.4.0-0.dev.0.20211022004519-f67a49cb50f5
replace github.com/maxlandon/readline => github.com/Rosettea/readline-1 v0.1.0-beta.0.20211207003625-341c7985ad7d
