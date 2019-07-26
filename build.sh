rm -rf cmd/konachan-app/static
cp -rf web/static cmd/konachan-app/
cd cmd/konachan-app
go build .