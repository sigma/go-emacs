[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-Ready--to--Code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/sigma/go-emacs) 

# goemacs
experimental Emacs modules in Go

Requires dynamic modules support to be built into Emacs (obviously): check that
`module-file-suffix` is not `nil`.

# Example

```
$ cd examples/hello
$ make test
emacs --batch --load hello.so --eval '(when (featurep (quote hello)) (hello "world"))'
INFO[0000] module initialization started
hello from go init
INFO[0000] creating native function
INFO[0000] creating symbol
INFO[0000] calling function
Hello function!
INFO[0000] calling symbol before it's bound
symbol is not a function
INFO[0000] binding symbol to function
INFO[0000] calling symbol after it's bound
Hello symbol!
INFO[0000] module initialization complete
Hello world!
```

# Credits
The following resources have been immensely useful:
* http://diobla.info/blog-archive/modules-tut.html
* http://nullprogram.com/blog/2016/11/05/
* https://blog.filippo.io/building-python-modules-with-go-1-5/
* http://blog.ralch.com/tutorial/golang-sharing-libraries/
* https://golang.org/cmd/cgo/
