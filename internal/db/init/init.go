package db_init

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/gudn/iinit"
	. "github.com/gudn/lesslog/internal/config"
	"github.com/gudn/lesslog/internal/config/init"
	"github.com/gudn/lesslog/internal/db"
)

func InitDb() {
	if C.Mode != "postgres" {
		return
	}
	pool, err := pgxpool.Connect(db.Ctx, C.Pg)
	if err != nil {
		log.Error().Err(err).Msg("failed connect to postgres")
	} else {
		log.Info().Msg("success connect to postgres")
		db.Pool = pool
	}
}

func init() {
	iinit.SequentialS(
		config_init.Init,
		InitDb,
	)
}
