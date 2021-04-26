BUILD		= imap-checker
OUTPUT		= $(BUILD)
VERSION		?= 0.0
PKG		?= $(BUILD)_$(VERSION)_amd64

all: clean format build

clean:
	rm -f $(OUTPUT)

format:
	go fmt

build:
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o $(OUTPUT) *.go

temp:
	$(eval TMPDIR = $(shell mktemp -d))
	$(eval OUTPUT = $(TMPDIR)/usr/bin/$(BUILD))
	mkdir -p $(TMPDIR)/etc/imap-checker
	mkdir -p $(TMPDIR)/usr/bin

files: temp
	cp -av ../DEBIAN $(TMPDIR)/
	cp -av ../config.json $(TMPDIR)/etc/imap-checker/

deb: clean format files build
	$(eval SIZE = $(shell du -k $(OUTPUT) | sed 's/^\([0-9]\+\).*/\1/'))
	$(shell sed -i "s/INSTALL_SIZE/$(SIZE)/" $(TMPDIR)/DEBIAN/control)
	$(shell sed -i "s/INSTALL_VERSION/$(VERSION)/" $(TMPDIR)/DEBIAN/control)
	dpkg-deb --build $(TMPDIR) $(PKG).deb
	rm -rf $(TMPDIR)
