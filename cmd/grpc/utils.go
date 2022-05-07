package main

import (
	"context"
	"errors"

	"github.com/gudn/lesslog/internal/logging"
	"google.golang.org/grpc"
)

func dismissContextCancel(
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
