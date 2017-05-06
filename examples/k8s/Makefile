GO ?= go
EMACS ?= emacs

k8s.so: main.go
	$(GO) build -buildmode=c-shared -ldflags -s -o $@ $<

clean:
	rm -f k8s.so

all: k8s.so

test: k8s.so
	$(EMACS) --batch --load k8s.so --eval '(message (format "%s" (k8s-list-pods (k8s-make-client (expand-file-name "~/.kube/config")))))'
