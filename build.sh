#!/bin/zsh
CGO_ENABLED=0 go build  -trimpath -a -ldflags="-w -s" -o ad2c .
