package service

import (
	"context"

	"github.com/gudn/lesslog/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func makeErr(name string) error {
	return status.Errorf(codes.Unimplemented, "method %v not implemented", name)
}

type UnimplementedService struct{}

func (UnimplementedService) Create(
	context.Context,
	string,
) (uint64, error) {
	return 0, makeErr("Create")
}

func (UnimplementedService) Push(
	context.Context,
	string,
	uint64,
	[]*proto.Operation,
) (bool, uint64, error) {
	return false, 0, makeErr("Push")
}

func (UnimplementedService) Fetch(
	context.Context,
	string,
	uint64,
	uint,
) ([]*proto.Operation, error) {
	return nil, makeErr("Fetch")
}

func (UnimplementedService) Watch(
	context.Context,
	string,
	uint64,
	uint,
) (<-chan []*proto.Operation, error) {
	return nil, makeErr("Watch")
}
