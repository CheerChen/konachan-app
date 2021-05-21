GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s'  -o
BIN=bin/konakore
SOURCE=.

docker:
	$(GOBUILD) $(BIN) $(SOURCE)