GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --tags --always --match "v*")
GIT_VERSION=github.com/vine-io/gpm/pkg/runtime/doc
CGO_ENABLED=0
BUILD_DATE=$(shell date +%s)
LDFLAGS=-X $(GIT_VERSION).GitCommit=$(GIT_COMMIT) -X $(GIT_VERSION).GitTag=$(GIT_TAG) -X $(GIT_VERSION).BuildDate=$(BUILD_DATE)

build-tag:
	sed -i "" "s/GitTag     = ".*"/GitTag     = \"$(GIT_TAG)\"/g" pkg/runtime/doc.go
	sed -i "" "s/GitCommit  = ".*"/GitCommit  = \"$(GIT_COMMIT)\"/g" pkg/runtime/doc.go
	sed -i "" "s/BuildDate  = ".*"/BuildDate  = \"$(BUILD_DATE)\"/g" pkg/runtime/doc.go

install:
	go mod vendor

build-darwin:
	mkdir -p cmd/gpm/pkg/testdata
	mkdir -p _output/darwin
	GOOS=darwin GOARCH=amd64 go build -o cmd/gpm/pkg/testdata/gpmd -a -installsuffix cgo -ldflags "-s -w" cmd/gpmd/main.go
	GOOS=darwin GOARCH=amd64 go build -o _output/darwin/gpm -a -installsuffix cgo -ldflags "-s -w" cmd/gpm/main.go

build-windows:
	mkdir -p cmd/gpm/pkg/testdata
	mkdir -p _output/windows
	cp nssm.exe cmd/gpm/pkg/testdata/nssm.exe
	GOOS=windows GOARCH=amd64 go build -o cmd/gpm/pkg/testdata/gpmd.exe -a -installsuffix cgo -ldflags "-s -w" cmd/gpmd/main.go
	GOOS=windows GOARCH=amd64 go build -o _output/windows/gpm.exe -a -installsuffix cgo -ldflags "-s -w" cmd/gpm/main.go

build-linux:
	mkdir -p cmd/gpm/pkg/testdata
	mkdir -p _output/linux
	GOOS=linux GOARCH=amd64 go build -o cmd/gpm/pkg/testdata/gpmd -a -installsuffix cgo -ldflags "-s -w" cmd/gpmd/main.go
	GOOS=linux GOARCH=amd64 go build -o _output/linux/gpm -a -installsuffix cgo -ldflags "-s -w" cmd/gpm/main.go

build: build-darwin build-linux build-windows

tar: build
	cd _output && \
	tar -zcvf gpm-darwin-$(GIT_TAG).tar.gz darwin/* && \
	tar -zcvf gpm-linux-$(GIT_TAG).tar.gz linux/*  && \
	zip gpm-windows-$(GIT_TAG).zip windows/* && \
	rm -fr darwin/ linux/ windows/

clean:
	rm -fr vendor

.PHONY: build-tag install build-darwin build-windows build-linux build tar clean