.PHONY: all clean

all: 
	gogll g.md
	go test -v
	dot -Tpng beforeDA.dot -o beforeDA.png
	dot -Tpng afterDA.dot -o afterDA.png

clean:
	rm *.txt; \
	rm -rf parser; \
	rm -rf goutil