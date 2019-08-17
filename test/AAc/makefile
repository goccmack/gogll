.Phony: clean test view

test:
	gogll aac.bnf && \
	go build && \
	./AAc && \
	dot -Tpng -osppf.png sppf.dot
	open sppf.png

view: 
	dot -Tpng -osppf.png sppf.dot
	open sppf.png

clean:
	rm *.txt; \
	rm -rf parser
