.PHONY: test clean

test:
	gogll -p gogll/test/test1 test1.bnf
	go test -v

clean:
	rm -rf goutil; \
	rm -rf parser; \
	rm *.txt