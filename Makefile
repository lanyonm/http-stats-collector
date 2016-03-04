.PHONY: test build help

BINARY := http-stats-collector

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

all: clean test build ## clean, test, and build all the things!

clean: ## delete dat binary!
	rm -f $(BINARY)

build: *.go ## mmmm, compilation!
	go build -o $(BINARY) *.go

test: ## avoid the footgun!
	go test -v

run: all ## equivalent of doit.sh
	./$(BINARY)

cover: ## get you some test coverage!
	go test -race -covermode=count -coverprofile=coverage.out && go tool cover -html=coverage.out
