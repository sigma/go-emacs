module github.com/sigma/go-emacs/examples/k8s

go 1.16

replace github.com/sigma/go-emacs => ../..

require (
	github.com/sigma/go-emacs v0.0.0-00010101000000-000000000000
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
)
