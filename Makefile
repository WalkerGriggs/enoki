.PHONY: all clean grpc-proto grpc-server grpc-proxy

all: grpc-proto grpc-server grpc-proxy

clean:
	rm ./internal/proto/manifests/**/*.go

grpc-proto:
	protoc -I ./internal/proto/ \
		-I ./third_party/googleapis/ \
		--go_out ./internal/proto \
		--go_opt paths=source_relative \
		internal/proto/manifests/v1/*.proto

grpc-server:
	protoc -I ./internal/proto/ \
		-I ./third_party/googleapis/ \
		--go-grpc_out ./internal/proto \
		--go-grpc_opt paths=source_relative \
		--go-grpc_opt require_unimplemented_servers=false \
		internal/proto/manifests/v1/*.proto

grpc-proxy:
	protoc -I ./internal/proto/ \
		-I ./third_party/googleapis/ \
		--grpc-gateway_out ./internal/proto \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		internal/proto/manifests/v1/*.proto
