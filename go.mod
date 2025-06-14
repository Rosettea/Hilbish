module hilbish

go 1.23.0

toolchain go1.24.3

require (
	github.com/aarzilli/golua v0.0.0-20250217091409-248753f411c4
	github.com/arnodel/golua v0.1.0
	github.com/atsushinee/go-markdown-generator v0.0.0-20231027094725-92d26ffbe778
	github.com/blackfireio/osinfo v1.1.0
	github.com/maxlandon/readline v1.1.3
	github.com/pborman/getopt v1.1.0
	github.com/sahilm/fuzzy v0.1.1
	golang.org/x/sys v0.33.0
	golang.org/x/term v0.32.0
	mvdan.cc/sh/v3 v3.11.0
)

require (
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/arnodel/strftime v0.1.6 // indirect
	github.com/evilsocket/islazy v1.11.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/olekukonko/ts v0.0.0-20171002115256-78ecb04241c0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/text v0.26.0 // indirect
)

replace mvdan.cc/sh/v3 => github.com/Rosettea/sh/v3 v3.4.0-0.dev.0.20240720131751-805c301321fd

replace github.com/maxlandon/readline => ./readline

replace layeh.com/gopher-luar => github.com/layeh/gopher-luar v1.0.10

replace github.com/arnodel/golua => github.com/Rosettea/golua v0.0.0-20241104031959-5551ea280f23
