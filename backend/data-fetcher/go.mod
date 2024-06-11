module tanken/backend/data-fetcher

go 1.22.2

require (
	connectrpc.com/connect v1.16.2
	github.com/aws/aws-sdk-go v1.53.1
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/redis/go-redis/v9 v9.5.1
	golang.org/x/net v0.25.0
	google.golang.org/protobuf v1.34.1
	tanken/backend/common v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/text v0.15.0 // indirect
)

replace tanken/backend/common => ../common
