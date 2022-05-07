package main

import (
	"context"

	. "github.com/gudn/lesslog/internal/config"
	"github.com/gudn/lesslog/pkg/service"
	"github.com/gudn/lesslog/pkg/service/pg"
	"github.com/gudn/lesslog/proto"
	"github.com/rs/zerolog/log"
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
	ops, err := l.s.Fetch(ctx, req.GetLogName(), req.GetSinceSn())
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
	ch, err := l.s.Watch(nested, req.GetLogName(), req.GetSinceSn())
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
	var s service.Interface
	switch C.Mode {
	case "unimplemented":
		log.Warn().Msg("use unimplemented service mode")
		s = service.UnimplementedService{}
	case "postgres":
		log.Info().Msg("user postgres service mode")
		s = pg.PostgresService{}
	default:
		log.Error().Str("mode", C.Mode).Msg("unrecognized mode; fallback to unimplemented")
		s = service.UnimplementedService{}
	}
	return &lesslogServer{s: s}
}
