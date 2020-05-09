SHELL:=/bin/bash

build:
	@ docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.25.1 golangci-lint run
	@ go version
	@ go clean
	@ env GO111MODULE=on GOOS="linux" GOARCH="$targetarch" go build -ldflags "-s -w" -o ehs main.go

build-release:
	@ chmod +x ../ci/build.sh
	@ ../ci/build.sh ${PWD}/../version.go "ehs"
