package main

import (
	"net"

	"github.com/gudn/iinit"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	. "github.com/gudn/lesslog/internal/config"
	_ "github.com/gudn/lesslog/internal/db/init"
	"github.com/gudn/lesslog/internal/logging"
	_ "github.com/gudn/lesslog/internal/metrics"
	"github.com/gudn/lesslog/proto"
)

func init() {
	iinit.Iinit()
}

func main() {
	lis, err := net.Listen("tcp", C.Bind)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	s := grpc.NewServer(
		grpc.StreamInterceptor(dismissContextCancel),
		grpc.UnaryInterceptor(logging.GrpcLogging),
	)
	proto.RegisterLesslogServer(s, Build())
	reflection.Register(s)

	log.Info().Str("bind", lis.Addr().String()).Msg("starting serving")
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
