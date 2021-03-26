# Hilbish
🎀 a nice lil shell for lua people made with go and lua

It is currently in a mostly beta state but is very much usable
(I'm using it right now).

# Links
- **[Documentation](https://github.com/Hilbis/Hilbish/wiki)**

# Building
Prebuilt binaries are not yet provided, so to try it out you'll have to manually compile.  

**NOTE:** Hilbish is currently only officially supported and tested on Linux

### Requirements
- Go 1.16

### Install
```sh
git clone https://github.com/Hilbis/Hilbish
cd Hilbish
make build
sudo make install
# Or 
sudo make
```

Alternativly, if you use Arch Linux, you can install Hilbish with an **(unofficial)** AUR package
```sh
yay -S hilbish-git
```

### Uninstall
```sh
sudo make uninstall
```

# License
[MIT](LICENSE)
