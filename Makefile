PREFIX ?= /usr
BINDIR ?= $(PREFIX)/bin
LIBDIR ?= $(PREFIX)/share/hilbish

build:
	@go build

install: build
	@install -v -d "$(BINDIR)/" && install -m 0755 -v hilbish "$(BINDIR)/hilbish"
	@mkdir -p "$(LIBDIR)"
	@cp libs preload.lua .hilbishrc.lua "$(LIBDIR)" -r
	@echo "Hilbish Installed"

uninstall:
	@rm -vrf \
			"$(BINDIR)/hilbish" \
			"$(LIBDIR)"
	@echo "Hilbish Uninstalled"

clean:
	@go clean

all: build install

.PHONY: install uninstall build clean
