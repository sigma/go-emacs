GO ?= go
EMACS ?= emacs

hello.so: main.go
	$(GO) build -buildmode=c-shared -o $@ $<

clean:
	rm -f hello.so

all: hello.so

test: hello.so
	$(EMACS) -batch -load hello.so
