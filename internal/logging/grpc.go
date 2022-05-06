package logging

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func logRequest(method string, err error) {
	if err != nil {
		log.
			Error().
			Err(err).
			Str("method", method).
			Msg("error processing request")
	} else {
		log.
			Info().
			Str("method", method).
			Msg("success processing request")
	}
}

func GrpcLogging(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	resp, err := handler(ctx, req)
	logRequest(info.FullMethod, err)
	return resp, err
}

func GrpcStreamLogging(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	err := handler(srv, ss)
	logRequest(info.FullMethod, err)
	return err
}
