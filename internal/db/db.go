package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Ctx = context.Background()
var Pool *pgxpool.Pool
