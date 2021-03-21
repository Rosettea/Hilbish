# Hilbish
🎀 a nice lil shell for lua people made with go and lua

It is currently in a mostly beta state but is very much usable
(I'm using it right now).

There are still some things missing like pipes but granted that will be
added soon.

# Links
- **[Documentation](https://github.com/Hilbis/Hilbish/wiki)**

# Building
Prebuilt binaries are not yet provided, so to try it out you'll have to manually compile.  

**NOTE:** Hilbish is currently only officially supported and tested on Linux

### Requirements
- Go 1.16

### Setup
```
git clone https://github.com/Hilbis/Hilbish
cd Hilbish
go build
```

This will build a `hilbish` executable in the current directory. 

# Install
`sudo cp hilbish /usr/bin`
`sudo mkdir /usr/share/hilbish`
`sudo cp libs preload.lua .hilbishrc.lua /usr/share/hilbish -r`

# License
[MIT](LICENSE)
