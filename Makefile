.PHONY: all grpc-start

all: proto goose-up

proto: lesslog.proto
	protoc --go_out=. --go-grpc_out=. lesslog.proto

goose-up: migrations/
	cd migrations && goose postgres "user=lesslog password=lesslog dbname=lesslog sslmode=disable" up

grpc-start:
	go run ./cmd/grpc -pretty -config config.dev
