.PHONY: all clean install

all: parser/parser.go
	go build

parser/parser.go: gogll.md
	gogll -p gogll gogll.md 

clean:
	first_follow.txt; \
	grammar_slots.txt; \
	symbols.txt; \
	rm -rf parser ; \
	rm -rf goutil/bsr ; \
	rm -rf goutil/md ; \
	rm -rf goutil/stringset ; \

