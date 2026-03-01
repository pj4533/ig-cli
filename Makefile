.PHONY: build test lint clean

build:
	go build -o ig-cli .

test:
	go test -race -coverprofile=cover.out ./...
	go-test-coverage --config=.testcoverage.yml

lint:
	golangci-lint run

clean:
	rm -f ig-cli cover.out
