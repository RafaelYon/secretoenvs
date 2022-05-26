# Go supported OSs and Archs: https://github.com/golang/go/blob/master/src/go/build/syslist.go

build: clear linux windows

linux: setup
	export GOOS=linux

	GOARCH=386 go build -o build/secretoenvs_linux_386
	GOARCH=amd64 go build -o build/secretoenvs_linux_amd64

	GOARCH=arm go build -o build/secretoenvs_linux_arm
	GOARCH=arm64 go build -o build/secretoenvs_linux_arm64

windows: setup
	export GOOS=windows

	GOARCH=386 go build -o build/secretoenvs_windows_386.exe
	GOARCH=amd64 go build -o build/secretoenvs_windows_amd64.exe
	
	GOARCH=arm go build -o build/secretoenvs_windows_arm.exe
	GOARCH=arm64 go build -o build/secretoenvs_windows_arm64.exe

setup:
	export CGO_ENABLED=0

clear:
	rm -rf build