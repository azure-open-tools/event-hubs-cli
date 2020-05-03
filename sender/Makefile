SHELL:=/bin/bash

build-local:
	@ docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.25.1 golangci-lint run
	@ go version
	@ go clean
	@ go build -o sender-bin main.go

build-release:
	@ docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.25.1 golangci-lint run
	@ chmod +x ../ci/build.sh
	@ ../ci/build.sh ${PWD}/../version.go "ehs"
