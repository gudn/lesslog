package logging

import (
	"flag"
	"os"

	"github.com/gudn/iinit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	. "github.com/gudn/lesslog/internal/config"
	"github.com/gudn/lesslog/internal/config/init"
)

func init() {
	pretty := flag.Bool("pretty", false, "")
	iinit.SequentialS(
		flag.Parse,
		func() {
			zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
			if *pretty {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			}
		},
	)

	iinit.SequentialS(
		config_init.Init,
		func() {
			level := C.Log.Level
			switch level {
			case "panic":
				zerolog.SetGlobalLevel(zerolog.PanicLevel)
			case "fatal":
				zerolog.SetGlobalLevel(zerolog.FatalLevel)
			case "error":
				zerolog.SetGlobalLevel(zerolog.ErrorLevel)
			case "warn":
				zerolog.SetGlobalLevel(zerolog.WarnLevel)
			case "info":
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			case "debug":
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			case "trace":
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			default:
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
		},
	)
}
