<div align="center">
	<img src="./assets/hilbish-flower.png" width=128><br>
	<img src="./assets/hilbish-text.png" width=256><br>
	<blockquote>
	🌺 The flower shell. A comfy and nice little shell for Lua users and fans!
	</blockquote>
	<h1></h1>
	<p align="center">
		<img alt="GitHub commit activity" src="https://img.shields.io/github/commit-activity/m/Rosettea/Hilbish?style=flat-square">
		<img alt="GitHub commits since latest release (by date)" src="https://img.shields.io/github/commits-since/Rosettea/Hilbish/latest?style=flat-square">
		<img alt="GitHub contributors" src="https://img.shields.io/github/contributors/Rosettea/Hilbish?style=flat-square"><br>
		<a href="https://github.com/Rosettea/Hilbish/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22"><img src="https://img.shields.io/github/issues/Hilbis/Hilbish/help%20wanted?style=flat-square&color=green" alt="help wanted"></a>
		<a href="https://github.com/Rosettea/Hilbish/blob/master/LICENSE"><img alt="GitHub license" src="https://img.shields.io/github/license/Rosettea/Hilbish?style=flat-square"></a>
		<a href="https://discord.gg/3PDdcQz"><img alt="Discord" src="https://img.shields.io/discord/732357621503229962?color=blue&style=flat-square"></a>
	</p>
</div>

Hilbish is a new, Unix-like shell which uses Lua as its configuration language.
It's decently fast for interactive use and includes a bunch of utility functions and features
to make your life in a terminal easier.

# Screenshots
<div align="center">
<img src="gallery/default.png"><br><br>
<img src="gallery/terminal.png"><br><br>
<img src="gallery/pillprompt.png">
</div>

# Installation
**NOTE:** Hilbish is currently only officially supported and tested on Linux

### Prebuilt binaries
Binaries are provided for the latest commit.  

**Note that these use Hilbiline, not readline, and may be missing functionality
(moving the cursor, proper unicode support and backspace working properly)**  

Click on the checkmark (or x) near the commit hash, then details for your platform  
<br><img src="https://modeus.is-inside.me/dyr8UGGq.png"><br>

Then click on the artifacts drop down, and download artifact for your platform,
like what is highlighted in the screenshot.  
<br><img src="https://modeus.is-inside.me/KJ0Puceb.png"><br>

### AUR
![AUR maintainer](https://img.shields.io/aur/maintainer/hilbish?logo=arch-linux&style=flat-square)  
Arch Linux users can install Hilbish from the AUR with the following command:  
```sh
yay -S hilbish
```
Or from the latest `master` commit with:  
```sh
yay -S hilbish-git
```

### Nixpkgs
Nix/NixOS users can install Hilbish from the central repository, nixpkgs, through the usual ways.

If you're new to nix you should probably read up on how to do that [here](https://nixos.wiki/wiki/Cheatsheet).

### Manual Build
#### Prerequisites
- [Go 1.16+](https://go.dev)
- GNU Readline

On Fedora or other RPM based distros, readline can be installed with:  
```
sudo dnf install readline-devel
```  

On Debian/Ubuntu and distros based on them, it can be installed with:  
```
sudo apt install libreadline-dev
```

On Arch Linux, it can be installed with:
```
sudo pacman -S readline
```

#### Build
First, clone Hilbish:
```sh
git clone --recursive https://github.com/Rosettea/Hilbish
cd Hilbish
go get -d ./...
```  

To build, run:
```
make dev
```  

Or, if you want a stable branch, run these commands:
```
git checkout $(git describe --tags `git rev-list --tags --max-count=1`)
make build
```  

After you did all that, run `sudo make install` to install Hilbish globally.

# Contributing
Any kind of contributions to Hilbish are welcome!
Read [CONTRIBUTING.md](CONTRIBUTING.md) before getting started.

##### Thanks to everyone below who's contributed!
<a href="https://github.com/Hilbis/Hilbish/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Hilbis/Hilbish" />
</a>

*Made with [contributors-img](https://contrib.rocks).*

# License
Hilbish is licensed under the MIT license!
Read the [license here](LICENSE) for more info. 

Images in the [assets](assets/) folder are under CC-BY-SA 4.0

