package main

import (
	"net"

	"github.com/gudn/iinit"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	. "github.com/gudn/lesslog/internal/config"
	_ "github.com/gudn/lesslog/internal/db/init"
	"github.com/gudn/lesslog/internal/metrics"
	_ "github.com/gudn/lesslog/internal/metrics"
	_ "github.com/gudn/lesslog/internal/mq/init"
	"github.com/gudn/lesslog/proto"
)

func init() {
	iinit.SequentialS(
		metrics.InitMetrics,
		func() {
			if metrics.IsEnabled() {
				prometheus.MustRegister(requests_total)
			}
		},
	)
	iinit.Iinit()
}

func main() {
	lis, err := net.Listen("tcp", C.Bind)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	s := grpc.NewServer(
		grpc.StreamInterceptor(streamMiddle),
		grpc.UnaryInterceptor(unaryMiddle),
	)
	proto.RegisterLesslogServer(s, Build())
	reflection.Register(s)

	log.Info().Str("bind", lis.Addr().String()).Msg("starting serving")
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
