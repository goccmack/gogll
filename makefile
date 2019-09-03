.PHONY: all clean install

all: parser/parser.go
	go install

parser/parser.go: gogll.md
	gogll -p gogll gogll.md 

clean:
	rm first_follow.txt; \
	rm grammar_slots.txt; \
	rm symbols.txt; \
	rm first_follow.txt; \
	rm -rf parser ; \
	rm -rf goutil/bsr ; \
	rm -rf goutil/md ; \
	rm -rf goutil/stringset ; \
	rm gogll

