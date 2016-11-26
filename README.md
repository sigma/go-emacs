# goemacs
experimental Emacs modules in Go

Requires dynamic modules support to be built into Emacs (obviously): check that
`module-file-suffix` is not `nil`.

# Example

```
$ cd examples/hello
$ make test
/usr/local/bin/emacs --batch --load hello.so --eval '(hello "world")'
hello from go init
Hello world!
```
