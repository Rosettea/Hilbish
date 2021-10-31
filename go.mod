module hilbish

go 1.16

require (
	github.com/Rosettea/Hilbiline v0.0.0-20210710124707-aa6e3ff34cb2
	github.com/bobappleyard/readline v0.0.0-20150707195538-7e300e02d38e
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pborman/getopt v1.1.0
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9
	golang.org/x/sys v0.0.0-20211030160813-b3129d9d1021 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	layeh.com/gopher-luar v1.0.10
	mvdan.cc/sh/v3 v3.4.0
)

replace mvdan.cc/sh/v3 => github.com/Rosettea/sh/v3 v3.4.0-0.dev.0.20211022004519-f67a49cb50f5
