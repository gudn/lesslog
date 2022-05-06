.PHONY: all grpc-start

all: proto

proto: lesslog.proto
	protoc --go_out=. --go-grpc_out=. lesslog.proto

grpc-start:
	go run ./cmd/grpc -pretty -config config.dev
