.PHONY: all

all: proto

proto: lesslog.proto
	protoc --go_out=. --go-grpc_out=. lesslog.proto
