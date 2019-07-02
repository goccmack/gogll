.PHONY: clean sppf

clean:
	rm -rf labels
	rm -rf sppf
	rm parser.go

sppf:
	go test -v
	dot -Tpng -osppf.png sppf.dot && \
	open sppf.png