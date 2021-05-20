GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'  -o
BIN=bin/konachan-app
SOURCE=.

docker:
	$(GOBUILD) $(BIN) $(SOURCE)
# linux-arm64:
# 	GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'  -o $(BIN)-linux-arm64 $(SOURCE)
# darwin-arm64:
# 	GOARCH=arm64 GOOS=linux CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'  -o $(BIN)-darwin-arm64 $(SOURCE)