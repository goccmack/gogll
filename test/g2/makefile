.PHONY: clean parser

parser:
	gogll -p gogll/test/g2 g2.bnf
	go test -v

clean: 
	rm -rf parser; \
	rm -rf util; \
	rm *.txt