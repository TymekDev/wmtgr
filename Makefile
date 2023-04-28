PREFIX = /usr/local
COMPLETIONS_DIR_BASH = $(PREFIX)/share/bash-completion/completions
COMPLETIONS_DIR_ZSH = $(PREFIX)/share/zsh/site-functions
COMPLETIONS_DIR_FISH = $(PREFIX)/share/fish/vendor_completions.d

all: wmtgr completions

wmtgr:
	 go build -ldflags "-X main.version=$$(git describe --always --dirty)" .

completions: wmtgr.bash wmtgr.zsh wmtgr.fish

wmtgr.bash: wmtgr
	./wmtgr --webmention-token '' completion bash > wmtgr.bash

wmtgr.zsh: wmtgr
	./wmtgr --webmention-token '' completion zsh > wmtgr.zsh

wmtgr.fish: wmtgr
	./wmtgr --webmention-token '' completion fish > wmtgr.fish

clean: 
	rm -f wmtgr wmtgr.bash wmtgr.zsh wmtgr.fish

install:
	install -d \
		$(PREFIX)/bin \
		$(COMPLETIONS_DIR_BASH) \
		$(COMPLETIONS_DIR_ZSH) \
		$(COMPLETIONS_DIR_FISH)

	install -pm 0755 wmtgr $(PREFIX)/bin/wmtgr
	install -pm 0644 wmtgr.bash $(COMPLETIONS_DIR_BASH)/wmtgr
	install -pm 0644 wmtgr.zsh $(COMPLETIONS_DIR_ZSH)/_wmtgr
	install -pm 0644 wmtgr.fish $(COMPLETIONS_DIR_FISH)/wmtgr.fish

uninstall:
	rm -f \
		$(PREFIX)/bin/wmtgr \
		$(COMPLETIONS_DIR_BASH)/wmtgr \
		$(COMPLETIONS_DIR_ZSH)/_wmtgr \
		$(COMPLETIONS_DIR_FISH)/wmtgr.fish

.PHONY: all wmtgr completions clean install uninstall
