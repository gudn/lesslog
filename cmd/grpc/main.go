package main

import (
	"github.com/gudn/iinit"
	"github.com/rs/zerolog/log"

	_ "github.com/gudn/lesslog/internal/config"
	_ "github.com/gudn/lesslog/internal/logging"
)

func init() {
	iinit.Iinit()
}

func main() {
	log.Info().Msg("hello world")
}
