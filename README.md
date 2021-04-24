<div align="center">
	<h1>Hilbish</h1>
	<blockquote>
	🎀 a nice lil shell for lua people made with go and lua
	</blockquote><p align="center">
		<a href="https://github.com/Hilbis/Hilbish/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22"><img src="https://img.shields.io/github/issues/Hilbis/Hilbish/help%20wanted?color=green" alt="help wanted"></a>
		<a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg"></a>
	</p>
</div>

Hilbish is an interactive Unix-like shell written in Go, with the config
and other code written in Lua.  
It is sort of in a stable state currently, usable as a daily shell,
but there may still be breaking changes in Lua modules.

# Screenshots
<div align="center">
<img src="gallery/default.png"><br><br>
<img src="gallery/terminal.png"><br><br>
<img src="gallery/pillprompt.png">
</div>

# Links
- **[Documentation](https://github.com/Hilbis/Hilbish/wiki)**
- **[Gallery](https://github.com/Hilbis/Hilbish/discussions/36)** - See
more screenshots of Hilbish in action

# Building
Prebuilt binaries are not yet provided, so to try it out you'll have to manually compile.  

**NOTE:** Hilbish is currently only officially supported and tested on Linux

### Prerequisites
- [Go 1.16](https://go.dev)

- GNU Readline

On Fedora, readline can be installed with:  
```
sudo dnf install readline-devel
```  

On Debian/Ubuntu and distros based on them, it can be installed with:  
```
sudo apt install libreadline-dev
```

### Install
```sh
git clone https://github.com/Hilbis/Hilbish
cd Hilbish
make build
sudo make install
# Or 
sudo make all
```

Alternativly, if you use Arch Linux, you can compile Hilbish with an **(unofficial)** AUR package
```sh
yay -S hilbish
```
If you want the latest and greatest, you can install and compile from latest git commit  
```sh
yay -S hilbish-git
```

### Uninstall
```sh
sudo make uninstall
```

# Contributing
Any kind of contributions to Hilbish are welcome!   
Make sure to read [CONTRIBUTING.md](CONTRIBUTING.md) before getting started.

### Special Thanks To
Everyone here who has contributed:
<a href="https://github.com/Hilbis/Hilbish/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Hilbis/Hilbish" />
</a>

*Made with [contributors-img](https://contrib.rocks).*

### Credits
- [This blog post](https://www.vidarholen.net/contents/blog/?p=878) which
is how Hilbish now inserts a newline even if output doesn't have one.

# License
Hilbish is licensed under the MIT license.  
[Read here](LICENSE) for more info.
