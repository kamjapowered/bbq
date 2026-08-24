version := "0.2.0"
go := "go"
cache := justfile_directory() / ".cache/go-build"

# everything ci runs, so a green local run means a green ci run
ci: fmt vet tidy race link

fmt:
	@test -z "$(gofmt -l .)" || { gofmt -l .; exit 1; }

vet:
	GOCACHE={{cache}} {{go}} vet ./...

tidy:
	{{go}} mod tidy -diff

test:
	GOCACHE={{cache}} {{go}} test ./...

race:
	GOCACHE={{cache}} {{go}} test -race ./...

# link every test binary with inlining off, see golang/go#80976
link:
	GOCACHE={{cache}} {{go}} test -gcflags="all=-N -l" -run='^$' ./...

microwave:
	PATH="$(GOTOOLCHAIN=local {{go}} env GOROOT)/bin:$PATH" microwave ./... --out bbq.go --pkg bbq

release: ci microwave
	git diff --check
	git status --short
	git add .
	git commit -m "v{{version}}"
	git tag -a v{{version}} -m "v{{version}}"
	git push --follow-tags
