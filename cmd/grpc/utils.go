package main

import (
	"context"
	"errors"

	"github.com/gudn/lesslog/internal/logging"
	"google.golang.org/grpc"
)

func streamMiddle(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	err := handler(srv, ss)
	if errors.Is(err, context.Canceled) {
		err = nil
	}
	logging.LogRequest(info.FullMethod, err)
	return err
}

func unaryMiddle(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	logging.LogRequest(info.FullMethod, err)
	return resp, err
}
