module github.com/gpm2/gpm

go 1.16

require (
	github.com/gogo/protobuf v1.3.2
	github.com/json-iterator/go v1.1.10
	github.com/lack-io/pkg/inject v0.2.0
	github.com/lack-io/plugins/dao/sqlite v0.18.0
	github.com/lack-io/vine v0.19.3
	google.golang.org/grpc v1.35.0
)

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
