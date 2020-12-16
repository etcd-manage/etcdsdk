module github.com/etcd-manage/etcdsdk

go 1.15

replace (
	github.com/coreos/bbolt => github.com/coreos/bbolt v1.3.3
	google.golang.org/grpc => google.golang.org/grpc v1.26.0 // grpc对etcd依赖问题
)

require (
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/prometheus/client_golang v1.8.0 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d // indirect
	google.golang.org/grpc v1.34.0 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
