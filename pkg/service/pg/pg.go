package pg

import (
	"context"
	"errors"

	"github.com/gudn/lesslog/internal/db"
	"github.com/gudn/lesslog/proto"
	"github.com/jackc/pgx/v4"
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
		return 0, nil
	}
	return *head_sn, nil
}

func (PostgresService) Push(
	ctx context.Context,
	log_name string,
	last_sn uint64,
	ops []*proto.Operation,
) (bool, uint64, error) {
	if db.Pool == nil {
		return false, 0, connIsNil
	}
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return false, 0, err
	}
	var head_sn *uint64
	err = tx.QueryRow(ctx, `
		SELECT head_sn
		FROM logs
		WHERE log_name = $1
    FOR UPDATE
	`, log_name).Scan(&head_sn)
	if errors.Is(err, pgx.ErrNoRows) {
		err = tx.Rollback(ctx)
		return false, 0, err
	} else if err != nil {
		_ = tx.Rollback(ctx)
		return false, 0, err
	}
	if (head_sn == nil && last_sn != 0) || (head_sn != nil && *head_sn != last_sn) {
		err = tx.Rollback(ctx)
		if head_sn == nil {
			return false, 0, err
		}
		return false, *head_sn, err
	}
	rows := make([][]any, len(ops))
	for i, op := range ops {
		last_sn += 1
		rows[i] = []any{log_name, last_sn, op.GetData()}
	}
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"operations"},
		[]string{"log_name", "sn", "data"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		_ = tx.Rollback(ctx)
		return false, 0, err
	}
	result, err := tx.Exec(ctx, `
		UPDATE logs SET head_sn = $1
		WHERE log_name = $2
	`, last_sn, log_name)
	if err != nil {
		_ = tx.Rollback(ctx)
		return false, 0, err
	}
	if result.RowsAffected() != 1 {
		err = tx.Rollback(ctx)
		return false, last_sn - uint64(len(ops)), err
	}
	err = tx.Commit(ctx)
	return true, last_sn, err
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
