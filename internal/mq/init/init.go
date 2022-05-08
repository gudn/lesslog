package mq_init

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/gudn/iinit"
	. "github.com/gudn/lesslog/internal/config"
	"github.com/gudn/lesslog/internal/logging"
	"github.com/gudn/lesslog/internal/mq"
)

func InitNats() {
	if C.Nats == "" {
		return
	}
	nc, err := nats.Connect(C.Nats)
	if err != nil {
		log.Error().Err(err).Msg("failed connect to nats")
	} else {
		log.Info().Msg("success connect to nats")
		mq.Conn = nc
	}
}

func init() {
	iinit.Sequential(
		logging.Init,
		iinit.Static(InitNats),
	)
}
