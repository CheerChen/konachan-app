#!/usr/bin/env bash
taskkill //im konachan-app.exe
go-bindata -o=internal/asset/asset.go -pkg=asset web/...
cd cmd/konachan-app
#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on GOPROXY=https://goproxy.io go build .
go build -ldflags '-w -s' .
#./upx.exe konachan-app