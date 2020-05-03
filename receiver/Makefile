SHELL:=/bin/bash

build-local:
	@ golangci-lint run
	@ go version
	@ go clean
	@ go build -o sender-bin sender/main.go

build-release:
	@ golangci-lint run
	@ chmod +x ../ci/build.sh
	@ ../ci/build.sh ${PWD}/../version.go "ehr"
