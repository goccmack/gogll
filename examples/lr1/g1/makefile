.PHONY: all go rust 

RUSTDIR = rust/g1
GODIR = go/src

all: go rust

go: $(GODIR)/go.mod
	../../../gogll -v -pager -o $(GODIR) g1.md

$(GODIR)/go.mod: 
	mkdir -p $(GODIR)
	cd $(GODIR); go mod init g1

rust: $(RUSTDIR)/Cargo.toml
	../../../gogll -a -v -t rust -pager -o $(RUSTDIR) g1.md

$(RUSTDIR)/Cargo.toml: 
	cargo new $(RUSTDIR)

