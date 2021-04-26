.PHONY: all clean install test

all: parser/parser.go
	go install

parser/parser.go: gogll.md
	gogll -o . gogll.md 

clean:
	rm first_follow.txt; \
	rm grammar_slots.txt; \
	rm symbols.txt; \
	rm -rf parser ; \
	rm gogll

test:
	make -C test

