taskkill //im konachan-app.exe
go-bindata -o=internal/asset/asset.go -pkg=asset web/...
cd cmd/konachan-app
go build -ldflags '-w -s' .
./upx.exe konachan-app.exe