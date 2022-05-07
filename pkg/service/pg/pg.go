package pg

import (
	"context"
	"errors"

	"github.com/gudn/lesslog/internal/db"
	"github.com/gudn/lesslog/pkg/messaging"
	"github.com/gudn/lesslog/proto"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var connIsNil error = status.Error(codes.Unavailable, "database is unavailable")
var messIsNil error = status.Error(codes.Unavailable, "messaging is unavailable")

type PostgresService struct {
	m messaging.Interface
}

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

func (p PostgresService) Push(
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
	if p.m != nil {
		p.m.Post(ctx, log_name, last_sn)
	}
	return true, last_sn, err
}

func (PostgresService) Fetch(
	ctx context.Context,
	log_name string,
	since_sn uint64,
	limit uint,
) ([]*proto.Operation, error) {
	if db.Pool == nil {
		return nil, connIsNil
	}
	rows, err := db.Pool.Query(ctx, `
		SELECT sn, data
		FROM operations
		WHERE log_name = $1 and sn > $2
		ORDER BY sn ASC
    LIMIT $3
  `, log_name, since_sn, limit)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	ops := make([]*proto.Operation, 0, limit)
	for rows.Next() {
		var sn uint64
		var data *[]byte
		if err := rows.Scan(&sn, &data); err != nil {
			return nil, err
		}
		op := &proto.Operation{Sn: sn}
		if data != nil {
			op.Data = *data
		}
		ops = append(ops, op)
	}
	return ops, nil
}

func (p PostgresService) Watch(
	ctx context.Context,
	log_name string,
	since_sn uint64,
	limit uint,
) (<-chan []*proto.Operation, error) {
	if db.Pool == nil {
		return nil, connIsNil
	}
	if p.m == nil {
		return nil, messIsNil
	}
	result := make(chan []*proto.Operation)

	var sendAll func() error
	sendAll = func() error {
		ops, err := p.Fetch(ctx, log_name, since_sn, limit)
		if err != nil {
			return err
		}
		if len(ops) > 0 {
			result <- ops
			since_sn = ops[len(ops)-1].Sn
			return sendAll()
		}
		return nil
	}

	ch, err := p.m.Listen(ctx, log_name)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := sendAll(); err != nil {
			log.Warn().Err(err).Msg("failed send new operations")
		}
		for sn := range ch {
			if sn > since_sn {
				err := sendAll()
				if err != nil {
					log.Warn().Err(err).Msg("failed send new operations")
					return
				}
			}
		}
	}()

	return result, nil
}


func New(m messaging.Interface) PostgresService {
	return PostgresService{m}
}
