.PHONY: gogll clean test

test: gogll
	make -C go test
	time cargo build --release
	time cargo run --release

gogll: 
	gogll -t rust gogll.md
	gogll -o go gogll.md

clean:
	rm -rf src/lexer; \
	rm -rf src/token; \
	rm -rf src/parser; \
	rm -rf target; \
	rm -rf go/lexer; \
	rm -rf go/parser; \
	rm -rf go/token 