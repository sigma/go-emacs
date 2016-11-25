# goemacs
experimental Emacs modules in Go

Requires dynamic modules support to be built into Emacs (obviously): check that
`module-file-suffix` is not `nil`.

# Example

```
$ cd examples/hello
$ make test
emacs -batch -load hello.so
hello from go
```
