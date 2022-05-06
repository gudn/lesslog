package config_init

import (
	"flag"

	"github.com/JeremyLoy/config"
	"github.com/gudn/iinit"
	"github.com/rs/zerolog/log"

	. "github.com/gudn/lesslog/internal/config"
)

var (
	configFile = flag.String("config", "config", "path to newline separated config file")
)

func Init() {
	err := config.FromOptional(*configFile).FromEnv().To(&C)
	if err != nil {
		log.Fatal().Err(err).Msg("configuration failed")
	}
}

func init() {
	iinit.SequentialS(
		flag.Parse,
		Init,
	)
}
