package metrics

import (
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/gudn/iinit"
	. "github.com/gudn/lesslog/internal/config"
	"github.com/gudn/lesslog/internal/logging"
)

var (
	mux = http.NewServeMux()
)

func InitMetrics() {
	if C.MonBind == "" {
		log.Warn().Msg("monitoring is disabled")
		return
	}
	lis, err := net.Listen("tcp", C.MonBind)
	if err != nil {
		log.Error().Err(err).Msg("failed to listen monitoring")
	}

	log.
		Info().
		Str("bind", lis.Addr().String()).
		Msg("starting serving monitoring")
	go func() {
		if err := http.Serve(lis, mux); err != nil {
			log.Error().Err(err).Msg("failed to serve monitoring")
		}
	}()
}

func init() {
	mux.Handle("/metrics", promhttp.Handler())
	iinit.Sequential(
		logging.Init,
		iinit.Static(InitMetrics),
	)
}
