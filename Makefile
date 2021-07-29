GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --tags --always --match "v*")
GIT_VERSION=github.com/gpm2/gpm/pkg/runtime/doc
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
	GOOS=darwin GOARCH=amd64 go build -o cmd/gpm/pkg/testdata/gpmd -a -installsuffix cgo -ldflags "-s -w" cmd/gpmd/main.go
	GOOS=darwin GOARCH=amd64 go build -o _output/gpm -a -installsuffix cgo -ldflags "-s -w" cmd/gpm/main.go

build-windows:

build-linux:
	mkdir -p cmd/gpm/pkg/testdata
	GOOS=darwin GOARCH=amd64 go build -o cmd/gpm/pkg/testdata/gpmd -a -installsuffix cgo -ldflags "-s -w" cmd/gpmd/main.go
	GOOS=darwin GOARCH=amd64 go build -o _output/gpm -a -installsuffix cgo -ldflags "-s -w" cmd/gpm/main.go

build: build-darwin

clean:
	rm -fr vendor

.PHONY: build-tag install build-darwin build-windows build-linux build clean