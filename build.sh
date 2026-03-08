#!/bin/bash

go build -o anmplugin-linux main.go

GOOS=darwin GOARCH=amd64 go build -o anmplugin-mac main.go
GOOS=darwin GOARCH=arm64 go build -o anmplugin-mac-arm main.go

GOOS=windows GOARCH=amd64 go build -o anmplugin.exe main.go