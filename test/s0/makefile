.PHONY: sppf

sppf:
	gogll s0.bnf && \
	go build && \
	./s0 && \
	dot -Tpng -osppf.png sppf.dot && \
	open sppf.png

view:
	dot -Tpng -osppf.png sppf.dot && \
	open sppf.png


clean:
	rm -rf parser ; \
	rm *.txt