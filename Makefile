.PHONY: test build

BINARY := http-stats-collector

all: clean test build

clean:
	rm -f $(BINARY)

build: *.go
	go build -o $(BINARY) *.go

test:
	go test -v ./...

run: all
	./$(BINARY)

cover:
	go test -race -covermode=count -coverprofile=coverage.out && go tool cover -html=coverage.out
