.DEFAULT_GOAL := build

-include .env

port=8080
token=$(if $(TOKEN),$(TOKEN), "")
name=dgpx
version=1.0

build: test
	go build

run: build
	echo "running application"
	go run *.go -port ${port} -token ${token}

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

build-image: build
	docker build -t "filipenos/${name}:${version}" .

run-image:
	touch .env
	docker run --rm --env-file=.env -p ${port}:8080 "filipenos/${name}:${version}"

push-image: build
	docker push "filipenos/${name}:${version}"