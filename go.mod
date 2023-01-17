module hilbish

go 1.17

require (
	github.com/arnodel/golua v0.0.0-20220221163911-dfcf252b6f86
	github.com/blackfireio/osinfo v1.0.3
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9
	github.com/pborman/getopt v1.1.0
	github.com/reeflective/readline v0.1.0-beta.0.20230108012254-bb995a76b5cd
	golang.org/x/sys v0.3.0
	golang.org/x/term v0.3.0
	mvdan.cc/sh/v3 v3.5.1
)

require (
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/arnodel/strftime v0.1.6 // indirect
	github.com/evilsocket/islazy v1.10.6 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/olekukonko/ts v0.0.0-20171002115256-78ecb04241c0 // indirect
	github.com/reiver/go-caret v0.0.0-20190726191135-0ad7aa63f55f // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace mvdan.cc/sh/v3 => github.com/Rosettea/sh/v3 v3.4.0-0.dev.0.20220524215627-dfd9a4fa219b

replace layeh.com/gopher-luar => github.com/layeh/gopher-luar v1.0.10

replace github.com/arnodel/golua => github.com/Rosettea/golua v0.0.0-20221213193027-cbf6d4e4d345

replace github.com/reeflective/readline => github.com/Rosettea/rl v0.0.0-20230114144324-1e2af14177c3
