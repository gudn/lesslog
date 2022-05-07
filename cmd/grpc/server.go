package main

import (
	"context"
	"strconv"

	. "github.com/gudn/lesslog/internal/config"
	"github.com/gudn/lesslog/pkg/messaging"
	"github.com/gudn/lesslog/pkg/messaging/local"
	"github.com/gudn/lesslog/pkg/service"
	"github.com/gudn/lesslog/pkg/service/pg"
	"github.com/gudn/lesslog/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
)

type lesslogServer struct {
	s service.Interface
	proto.UnimplementedLesslogServer
}

func (l *lesslogServer) Create(
	ctx context.Context,
	req *proto.CreateRequest,
) (*proto.PushResponse, error) {
	last_sn, err := l.s.Create(ctx, req.GetLogName())
	if err != nil {
		return nil, err
	}
	return &proto.PushResponse{Success: true, LastSn: last_sn}, nil
}

func (l *lesslogServer) Push(
	ctx context.Context,
	req *proto.PushRequest,
) (*proto.PushResponse, error) {
	succ, last_sn, err := l.s.Push(
		ctx,
		req.GetLogName(),
		req.GetLastSn(),
		req.GetOperations(),
	)
	if err != nil {
		return nil, err
	}
	return &proto.PushResponse{Success: succ, LastSn: last_sn}, nil
}

func (l *lesslogServer) Fetch(
	ctx context.Context,
	req *proto.FetchRequest,
) (*proto.FetchResponse, error) {
	limit := uint(5)
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		val := md["x-limit"]
		if len(val) >= 1 {
			val, err := strconv.ParseUint(val[0], 10, 32)
			if err != nil {
				return nil, err
			}
			limit = uint(val)
		}
	}
	ops, err := l.s.Fetch(ctx, req.GetLogName(), req.GetSinceSn(), limit)
	if err != nil {
		return nil, err
	}
	return &proto.FetchResponse{Operations: ops}, nil
}

func (l *lesslogServer) Watch(
	req *proto.FetchRequest,
	stream proto.Lesslog_WatchServer,
) error {
	ctx := stream.Context()
	nested, cancel := context.WithCancel(ctx)
	defer cancel()
	ch, err := l.s.Watch(nested, req.GetLogName(), req.GetSinceSn(), 5)
	if err != nil {
		return err
	}
	for {
		select {
		case ops := <-ch:
			if err := stream.Send(&proto.FetchResponse{Operations: ops}); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func Build() *lesslogServer {
	var m messaging.Interface
	switch C.Messaging {
	case "local":
		log.Warn().Msg("user local messaging")
		m = local.New()
	case "none":
		log.Warn().Msg("disable messaging")
	default:
		log.Error().Str("messaging", C.Messaging).Msg("unrecognized messaging mode; disable messaging")
	}
	var s service.Interface
	switch C.Mode {
	case "unimplemented":
		log.Warn().Msg("use unimplemented service mode")
		s = service.UnimplementedService{}
	case "postgres":
		log.Info().Msg("user postgres service mode")
		s = pg.New(m)
	default:
		log.Error().Str("mode", C.Mode).Msg("unrecognized mode; fallback to unimplemented")
		s = service.UnimplementedService{}
	}
	return &lesslogServer{s: s}
}
