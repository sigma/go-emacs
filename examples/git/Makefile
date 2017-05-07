GO ?= go
EMACS ?= emacs

git.so: main.go
	$(GO) build -buildmode=c-shared -ldflags -s -o $@ $<

clean:
	rm -f git.so

all: git.so

test: git.so
	$(EMACS) --batch --load $< --eval '(message "%s" (when (featurep (quote git)) (let ((repo (git-open (expand-file-name default-directory)))) (git-ls-branches repo))))'
