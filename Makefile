SHELL:=/bin/bash

build-local-sender:
	go version
	go clean
	go build -ldflags "-s -w" -o sender/bin/ehst sender/main.go

build-local-receiver:
	go version
	go clean
	go build -ldflags "-s -w" -o receiver/bin/ehrt receiver/main.go

build-sender:
	@go version
	@echo $(OS)
ifeq ($(OS),Windows_NT)
	set GOARCH=$(targetarch)
	set GOOS=$(targetos)
	go build -ldflags "-s -w" -o ehst-$(buildVersion)$(extension) sender/main.go
	cp ehst-$(buildVersion)$(extension) $(output)
	ls -lah $(output)
else
	env GOOS=$(targetos) GOARCH=$(targetarch) go build -ldflags "-s -w" -o $(output)/ehst-$(buildVersion)$(extension) sender/main.go
endif

build-receiver:
	@go version
ifeq ($(OS),Windows_NT)
	set GOARCH=$(targetarch)
	set GOOS=$(targetos)
	go build -ldflags "-s -w" -o ehrt-$(buildVersion)$(extension) receiver/main.go
	cp ehrt-$(buildVersion)$(extension) $(output)
	ls -lah $(output)
else
	env GOOS=$(targetos) GOARCH=$(targetarch) go build -ldflags "-s -w" -o $(output)/ehrt-$(buildVersion)$(extension) receiver/main.go
endif
