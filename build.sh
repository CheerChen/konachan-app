#!/usr/bin/env bash
cd cmd/konachan-app
GOARCH=arm64 GOOS=linux go build -trimpath -ldflags '-w -s' .
# ./konachan-app.exe