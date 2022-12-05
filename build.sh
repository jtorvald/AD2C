#!/bin/zsh
CGO_ENABLED=0 go build  -trimpath -a -ldflags="-w -s" -o ad2c .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build  -trimpath -a -ldflags="-w -s" -o ad2c.exe .
