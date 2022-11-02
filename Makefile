# Go supported OSs and Archs: https://github.com/golang/go/blob/master/src/go/build/syslist.go

export CGO_ENABLED=0

build: clear linux windows sha256results

linux:
	GOOS=linux GOARCH=386 go build -o build/secretoenvs_linux_386
	GOOS=linux GOARCH=amd64 go build -o build/secretoenvs_linux_amd64

	GOOS=linux GOARCH=arm go build -o build/secretoenvs_linux_arm
	GOOS=linux GOARCH=arm64 go build -o build/secretoenvs_linux_arm64

windows:
	GOOS=windows GOARCH=386 go build -o build/secretoenvs_windows_386.exe
	GOOS=windows GOARCH=amd64 go build -o build/secretoenvs_windows_amd64.exe

	GOOS=windows GOARCH=arm go build -o build/secretoenvs_windows_arm.exe
	GOOS=windows GOARCH=arm64 go build -o build/secretoenvs_windows_arm64.exe

sha256results:
	cd build && sha256sum * > sha256sums.txt

clear:
	rm -rf build