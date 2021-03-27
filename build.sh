#!/usr/bin/env bash
cd cmd/konachan-app
go build -trimpath -ldflags '-w -s' .
./konachan-app.exe