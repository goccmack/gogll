.PHONY: all clean

all:
	gogll -p gogll/test/s1 s1.bnf && \
	go test -v

clean:
	rm -rf goutil; \
	rm -rf parser; \
	rm *.txt