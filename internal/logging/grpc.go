package logging

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func LogRequest(method string, err error) {
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
	LogRequest(info.FullMethod, err)
	return resp, err
}
