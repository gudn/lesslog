package main

import (
	"context"
	"flag"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/gudn/iinit"
	_ "github.com/gudn/lesslog/internal/logging"
	"github.com/gudn/lesslog/proto"
)

var (
	url      = flag.String("url", "localhost:8080", "url to lesslog server")
	log_name = flag.String("log", "log", "log name")
	since    = flag.Uint64("since", 0, "since argument to watch")
)

func init() {
	iinit.Static(flag.Parse)
	iinit.Iinit()
}

func main() {
	ctx := context.Background()

	dial, err := grpc.DialContext(ctx, *url, grpc.WithInsecure())
	if err != nil {
		log.
			Fatal().
			Err(err).
			Str("url", *url).
			Msg("failed connect")
	}
	defer dial.Close()

	client := proto.NewLesslogClient(dial)
	ss, err := client.Watch(ctx, &proto.FetchRequest{LogName: *log_name, SinceSn: *since})
	if err != nil {
		log.
			Fatal().
			Err(err).
			Msg("failed execute method")
	}

	for {
		ops, err := ss.Recv()
		if err != nil {
			log.
				Fatal().
				Err(err).
				Msg("failed recv operations")
		}
		for _, op := range ops.Operations {
			log.
				Info().
				Uint64("sn", op.GetSn()).
				Bytes("data", op.GetData()).
				Msg("received operation")
		}
	}
}
