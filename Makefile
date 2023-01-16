.PHONY: all clean grpc-proto grpc-server grpc-proxy

all: grpc-proto grpc-server grpc-proxy

clean:
	rm ./proto/manifests/**/*.go

grpc-proto:
	protoc -I ./proto/ \
		--go_out ./proto \
		--go_opt paths=source_relative \
		proto/manifests/v1/*.proto

grpc-server:
	protoc -I ./proto/ \
		--go-grpc_out ./proto \
		--go-grpc_opt paths=source_relative \
		--go-grpc_opt require_unimplemented_servers=false \
		proto/manifests/v1/*.proto

grpc-proxy:
	protoc -I ./proto/ --grpc-gateway_out ./proto \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		proto/manifests/v1/*.proto
