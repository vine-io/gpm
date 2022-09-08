GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --tags --always --match "v*")
GIT_VERSION=github.com/vine-io/gpm/pkg/runtime/doc
CGO_ENABLED=0
BUILD_DATE=$(shell date +%s)
LDFLAGS=-X $(GIT_VERSION).GitCommit=$(GIT_COMMIT) -X $(GIT_VERSION).GitTag=$(GIT_TAG) -X $(GIT_VERSION).BuildDate=$(BUILD_DATE)

release:
ifeq "$(TAG)" ""
	@echo "missing tag"
	exit 1
endif
	git tag $(TAG)
	make build-tag
	git add .
	git commit -m "$(TAG)"
	git tag -d $(TAG)
	git tag $(TAG)

build-tag:
	sed -i "" "s/GitTag     = ".*"/GitTag     = \"$(GIT_TAG)\"/g" pkg/runtime/doc.go
	sed -i "" "s/GitCommit  = ".*"/GitCommit  = \"$(GIT_COMMIT)\"/g" pkg/runtime/doc.go
	sed -i "" "s/BuildDate  = ".*"/BuildDate  = \"$(BUILD_DATE)\"/g" pkg/runtime/doc.go

install:
	go install github.com/vine-io/vine/cmd/vine
	go install github.com/vine-io/vine/cmd/protoc-gen-gogo
	go install github.com/vine-io/vine/cmd/protoc-gen-vine
	go install github.com/vine-io/vine/cmd/protoc-gen-deepcopy
	go install github.com/vine-io/vine/cmd/protoc-gen-validator
	go install github.com/vine-io/vine/cmd/protoc-gen-dao

generate-ssl:
	cd build/ssl && \
	cfssl gencert -initca ca-csr.json | cfssljson -bare ca - && \
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=localhost server-csr.json | cfssljson -bare server && \
	cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=localhost client-csr.json | cfssljson -bare client && \
  	/bin/mv ca.pem ../../pkg/runtime/ssl && /bin/mv ca-key.pem ../../pkg/runtime/ssl && \
  	/bin/mv client.pem ../../pkg/runtime/ssl && /bin/mv client-key.pem ../../pkg/runtime/ssl && \
  	/bin/mv server.pem ../../pkg/runtime/ssl && /bin/mv server-key.pem ../../pkg/runtime/ssl


vendor:
	go mod vendor

build-darwin-amd64:
	mkdir -p cmd/gpm/pkg/testdata
	mkdir -p _output/darwin
	GOOS=darwin GOARCH=amd64 go build -o cmd/gpm/pkg/testdata/gpmd -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" cmd/gpmd/main.go
	GOOS=darwin GOARCH=amd64 go build -o _output/darwin/gpm -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" cmd/gpm/main.go

build-darwin-arm64:
	mkdir -p cmd/gpm/pkg/testdata
	mkdir -p _output/darwin
	GOOS=darwin GOARCH=arm64 go build -o cmd/gpm/pkg/testdata/gpmd -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" cmd/gpmd/main.go
	GOOS=darwin GOARCH=arm64 go build -o _output/darwin/gpm -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" cmd/gpm/main.go

build-windows:
	mkdir -p cmd/gpm/pkg/testdata
	mkdir -p _output/windows
	cp nssm.exe cmd/gpm/pkg/testdata/nssm.exe
	GOOS=windows GOARCH=amd64 go build -o cmd/gpm/pkg/testdata/gpmd.exe -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" cmd/gpmd/main.go
	GOOS=windows GOARCH=amd64 go build -o _output/windows/gpm.exe -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" cmd/gpm/main.go

build-linux-amd64:
	mkdir -p cmd/gpm/pkg/testdata
	mkdir -p _output/linux
	GOOS=linux GOARCH=amd64 go build -o cmd/gpm/pkg/testdata/gpmd -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" cmd/gpmd/main.go
	GOOS=linux GOARCH=amd64 go build -o _output/linux/gpm -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" cmd/gpm/main.go

build-linux-arm64:
	mkdir -p cmd/gpm/pkg/testdata
	mkdir -p _output/linux
	GOOS=linux GOARCH=arm64 go build -o cmd/gpm/pkg/testdata/gpmd -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" cmd/gpmd/main.go
	GOOS=linux GOARCH=arm64 go build -o _output/linux/gpm -a -installsuffix cgo -ldflags "-s -w ${LDFLAGS}" cmd/gpm/main.go

build-amd: build-darwin-amd64 build-linux-amd64 build-windows

build-arm: build-darwin-arm64 build-linux-arm64 build-windows

build: build-amd build-arm

tar-amd: build-amd
	cd _output && \
	tar -zcvf gpm-darwin-amd64-$(GIT_TAG).tar.gz darwin/* && \
	tar -zcvf gpm-linux-amd64-$(GIT_TAG).tar.gz linux/*  && \
	zip gpm-windows-$(GIT_TAG).zip windows/* && \
	rm -fr darwin/ linux/ windows/

tar-arm: build-arm
	cd _output && \
	tar -zcvf gpm-darwin-arm64-$(GIT_TAG).tar.gz darwin/* && \
	tar -zcvf gpm-linux-arm64-$(GIT_TAG).tar.gz linux/*  && \
	rm -fr darwin/ linux/

tar: generate-ssl tar-amd tar-arm

test-coverage:
	go test ./... -bench=. -coverage

lint:
	golint .

clean:
	rm -fr vendor

.PHONY: release build-tag vendor install build-darwin build-windows build-linux build tar clean