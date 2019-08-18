.PHONY: all clean install

all: parser/parser.go
	go build

parser/parser.go: gogll.md
	gogll -p gogll gogll.md 

clean:
	rm *..txt ; \
	rm -rf parser ; \

