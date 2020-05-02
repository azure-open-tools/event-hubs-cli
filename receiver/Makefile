SHELL:=/bin/bash

lint:
	@ docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.25.1 golangci-lint run

build-local: lint
	@ go version
	@ go clean
	@ go build -o sender-bin sender/main.go

build-release: lint
	@ chmod +x ../ci/build.sh
	@ ../ci/build.sh ${PWD}/../version.go "ehr"
