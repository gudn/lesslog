package db_init

import (
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/gudn/iinit"
	. "github.com/gudn/lesslog/internal/config"
	"github.com/gudn/lesslog/internal/db"
	"github.com/gudn/lesslog/internal/logging"
)

func InitDb() {
	if C.Mode != "postgres" {
		return
	}
	conf, err := pgxpool.ParseConfig(C.Pg)
	if err != nil {
		log.Error().Err(err).Msg("failed parse to Pg config")
		return
	}
	// NOTE postgres log with INFO level
	if log.Debug().Enabled() {
		conf.ConnConfig.Logger = zerologadapter.NewLogger(log.Logger)
		conf.ConnConfig.LogLevel, err = pgx.LogLevelFromString(C.Log.Level)
		if err != nil {
			conf.ConnConfig.LogLevel = pgx.LogLevelInfo
		}
	}
	pool, err := pgxpool.ConnectConfig(db.Ctx, conf)
	if err != nil {
		log.Error().Err(err).Msg("failed connect to postgres")
	} else {
		log.Info().Msg("success connect to postgres")
		db.Pool = pool
	}
}

func init() {
	iinit.Sequential(
		logging.Init,
		iinit.Static(InitDb),
	)
}
