module github.com/gpm2/gpm

go 1.16

require (
	github.com/StackExchange/wmi v1.2.0 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/json-iterator/go v1.1.10
	github.com/lack-io/cli v1.2.4
	github.com/lack-io/pkg/inject v0.2.0
	github.com/lack-io/plugins/dao/postgres v0.18.0
	github.com/lack-io/plugins/dao/sqlite v0.18.0
	github.com/lack-io/vine v0.19.3
	github.com/shirou/gopsutil v3.21.6+incompatible
	github.com/tklauser/go-sysconf v0.3.7 // indirect
	google.golang.org/grpc v1.35.0
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
