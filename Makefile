PREFIX = /usr/local
COMPLETIONS_DIR_FISH = $(PREFIX)/share/fish/vendor_completions.d

all: wmtgr completions

wmtgr:
	 go build -ldflags "-X main.version=$$(git describe --always --dirty)" .

completions: wmtgr.fish

wmtgr.fish: wmtgr
	./wmtgr --webmention-token '' completion fish > wmtgr.fish

clean: 
	rm -f wmtgr wmtgr.fish

install:
	install -d \
		$(PREFIX)/bin \
		$(COMPLETIONS_DIR_FISH)

	install -pm 0755 wmtgr $(PREFIX)/bin/wmtgr
	install -pm 0644 wmtgr.fish $(COMPLETIONS_DIR_FISH)/wmtgr.fish

uninstall:
	rm -f \
		$(PREFIX)/bin/wmtgr \
		$(COMPLETIONS_DIR_FISH)/wmtgr.fish

.PHONY: all wmtgr completions clean install uninstall
