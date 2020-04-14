.PHONY: all clean install

all: parser/parser.go
	go install

parser/parser.go: gogll.md
	gocc -p github.com/goccmack/gogll gogll.md 

clean:
	rm first_follow.txt; \
	rm grammar_slots.txt; \
	rm symbols.txt; \
	rm -rf parser ; \
	rm gogll

