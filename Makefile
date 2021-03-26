PREFIX ?= /usr
BINDIR ?= $(PREFIX)/bin
LIBDIR ?= $(PREFIX)/share/hilbish

build:
	@go build

install:
	@install -v -d "$(BINDIR)/" && install -m 0755 -v hilbish "$(BINDIR)/hilbish"
	@mkdir -p "$(LIBDIR)"
	@cp libs preload.lua .hilbishrc.lua "$(LIBDIR)" -r
	@echo /usr/bin/hilbish >> /etc/shells
	@echo "Hilbish Installed"

uninstall:
	@rm -vrf \
			"$(BINDIR)/hilbish" \
			"$(LIBDIR)"
	@sed '/\/usr\/bin\/hilbish/d' /etc/shells
	@echo "Hilbish Uninstalled"

clean:
	@go clean

all: build install

.PHONY: install uninstall build clean
