version := "0.3.0"
go := "go"
cache := justfile_directory() / ".cache/go-build"

# everything ci runs, so a green local run means a green ci run. gowork=off so
# a sibling in the workspace cannot shadow a published dep and hide a build
# that only works on this machine
ci:
	GOWORK=off just fmt vet tidy race link

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
