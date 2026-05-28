version := "0.1.0"

test:
	GOCACHE={{justfile_directory()}}/.cache/go-build go test ./...

release: test
	git diff --check
	git status --short
	git add .
	git commit -m "v{{version}}"
	git tag -a v{{version}} -m "v{{version}}"
	git push --follow-tags
