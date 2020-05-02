#!/usr/bin/env bash

versionFile=$1
name=$2

go version
echo "$OS"
echo "$targetarch"

targetos="$OS"
targetarch="amd64"
version=$(go run "$versionFile")

if [[ "$targetos" == *"Windows_NT"* ]];
then
	set GOARCH="$targetarch"
	set GOOS="$targetos"
	extension=".exe"
	go build -ldflags "-s -w" -o "$name-windows""$targetarch"-"$version""$extension" main.go
	mv "$name-windows""$targetarch"-"$version""$extension" ../
	ls -lah
else
  targetos=$(sw_vers | awk '{print $2$3$4}' | head -n 1)
  echo "Target OS: $targetos"
  if [[ "$targetos" == *"MacOSX"* ]];
  then
  	echo "$PWD"
    env GOOS="darwin" GOARCH="$targetarch" go build -ldflags "-s -w" -o "$name-darwin-""$targetarch"-"$version" main.go
  	mv "$name-darwin-""$targetarch"-"$version" ../
  else
  	env GOOS="linux" GOARCH="$targetarch" go build -ldflags "-s -w" -o "$name-linux-""$targetarch"-"$version" main.go
    mv "$name-linux-""$targetarch"-"$version" ../
  fi
fi