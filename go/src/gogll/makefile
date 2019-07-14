.PHONY: clean lexer gocc debug install

install: clean gocc
	go install

clean:
	rm first.txt ; \
	rm lexer_sets.txt ; \
	rm LR1_* ; \
	rm terminals.txt ; \
	rm -rf lexer ; \
	rm -rf parser ; \

gocc:
	gocc -p gogll gogll.bnf 

debug:
	gocc -p gogll -v -debug_lexer gogll.bnf