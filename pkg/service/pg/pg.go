package pg

import (
	"context"

	"github.com/gudn/lesslog/internal/db"
	"github.com/gudn/lesslog/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var connIsNil error = status.Error(codes.Unavailable, "database is unavailable")

type PostgresService struct{}

func (PostgresService) Create(
	ctx context.Context,
	log_name string,
) (uint64, error) {
	if db.Pool == nil {
		return 0, connIsNil
	}
	sql := `
		INSERT INTO logs(log_name)
		VALUES ($1)
		ON CONFLICT (log_name) DO UPDATE
		SET log_name = EXCLUDED.log_name
		RETURNING head_sn
  `
	var head_sn *uint64
	err := db.Pool.QueryRow(ctx, sql, log_name).Scan(&head_sn)
	if err != nil {
		return 0, err
	}
	if head_sn == nil {
		return uint64(0), nil
	}
	return *head_sn, nil
}

func (PostgresService) Push(
	context.Context,
	string,
	uint64,
	[]*proto.Operation,
) (bool, uint64, error) {
	if db.Pool == nil {
		return false, 0, connIsNil
	}
	panic("unimplemented")
}

func (PostgresService) Fetch(
	context.Context,
	string,
	uint64,
) ([]*proto.Operation, error) {
	if db.Pool == nil {
		return nil, connIsNil
	}
	panic("unimplemented")
}

func (PostgresService) Watch(
	context.Context,
	string,
	uint64,
) (<-chan []*proto.Operation, error) {
	if db.Pool == nil {
		return nil, connIsNil
	}
	panic("unimplemented")
}
