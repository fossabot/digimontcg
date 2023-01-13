.POSIX:
.SUFFIXES:

.PHONY: default
default: build

.PHONY: build
build:
	go build -o digimontcg main.go

.PHONY: web
web:
	go run main.go

.PHONY: update
update:
	go get -u ./...
	go mod tidy

.PHONY: test
test:
	go test -count=1 ./...

.PHONY: cover
cover:
	go test -coverprofile=c.out -coverpkg=./... -count=1 ./...
	go tool cover -html=c.out

.PHONY: release
release:
	goreleaser release --snapshot --rm-dist

.PHONY: format
format:
	gofmt -l -s -w .

.PHONY: clean
clean:
	rm -fr digimontcg c.out dist/
