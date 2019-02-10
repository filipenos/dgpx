.DEFAULT_GOAL := build

build: test
	go build ./...

run: build
	echo "running application"
	go run *.go

clean:
	go clean
	rm -f coverage.out

vet:
	go vet ./...

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out