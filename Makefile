BINARY := fetch

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build -o $(BINARY) main.go