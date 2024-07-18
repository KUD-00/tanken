module tanken/backend/user

go 1.22.5

require (
	connectrpc.com/connect v1.16.2
	github.com/aws/aws-sdk-go v1.54.19
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
	github.com/redis/go-redis/v9 v9.5.4
	golang.org/x/net v0.27.0
	google.golang.org/protobuf v1.33.0
	tanken/backend/common v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

replace tanken/backend/common => ../common
