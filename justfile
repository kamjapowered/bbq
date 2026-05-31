version := "0.1.1"

test:
	GOCACHE={{justfile_directory()}}/.cache/go-build go test ./...

microwave:
	microwave ./... --out bbq.go --pkg bbq

release: test microwave
	git diff --check
	git status --short
	git add .
	git commit -m "v{{version}}"
	git tag -a v{{version}} -m "v{{version}}"
	git push --follow-tags
