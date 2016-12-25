# k8s emacs plugin

Demo code for listing Kubernetes pods in Emacs, through https://github.com/kubernetes/client-go

# Usage

```
$ make test
emacs --batch --load k8s.so --eval '(message (format "%s" (k8s-list-pods (k8s-make-client (expand-file-name "~/.kube/config")))))'
(pod1 pod2)
```

# Notes

This is a little bit interesting, since `k8s-make-client` generates a client
object on the go side, which is passed as a user pointer to emacs. Then emacs
passes it back to the go code, so that the implementation of `k8s-list-pods`
can make use of it.

All in all, a lot of back and forth between Go and Emacs, related to
non-trivial objects.
