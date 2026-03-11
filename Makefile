whoisx: go.* *.go
	go build -v

.PHONY: clean
clean:
	rm -rf dist

.PHONY: lint
lint: whoisx
	go tool golangci-lint run

.PHONY: lint-fix
lint-fix: whoisx
	go tool golangci-lint run --fix

.PHONY: pre-commit
pre-commit: whoisx
	go tool golangci-lint run --disable godox

.PHONY: release
release:
	goreleaser release --snapshot --clean

