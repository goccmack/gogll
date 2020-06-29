.PHONY: all go rust clean

RUSTDIR = rust/g1
GODIR = go/src

all: go rust

go: $(GODIR)/go.mod
	gogll -pager -o $(GODIR) g1.md

rust: $(RUSTDIR)/Cargo.toml
	gogll -rust -pager -o $(RUSTDIR) g1.md

$(GODIR)/go.mod: 
	mkdir -p $(GODIR)
	cd $(GODIR); go mod init g1

$(RUSTDIR)/Cargo.toml: 
	cargo new $(RUSTDIR)

clean:
	rm go/src/*.txt; \
	rm $(RUSTDIR)/*.txt; \
	rm -rf $(RUSTDIR)/target; \
	rm -rf $(RUSTDIR)/src/errors; \
	rm -rf $(RUSTDIR)/src/lexer; \
	rm -rf $(RUSTDIR)/src/parser; \
	rm -rf $(RUSTDIR)/src/token; \
