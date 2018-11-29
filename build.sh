env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o builds/markpdf_darwin-amd64
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/markpdf_linux-amd64
env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o builds/markpdf_linux-arm64
env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o builds/markpdf_windows-386.exe
upx builds/markpdf_*