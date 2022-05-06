package service

import (
	"context"

	"github.com/gudn/lesslog/proto"
)

type Interface interface {
	Create(ctx context.Context, log_name string) (
		last_sn uint64,
		err error,
	)

	Push(
		ctx context.Context,
		log_name string,
		last_sn uint64,
		ops []*proto.Operation,
	) (
		success bool,
		last_sn_ uint64,
		err error,
	)

	Fetch(
		ctx context.Context,
		log_name string,
		since_sn uint64,
	) (
		ops []*proto.Operation,
		err error,
	)

	Watch(
		ctx context.Context,
		log_name string,
		since_sn uint64,
	) (
		ch <-chan []*proto.Operation,
		err error,
	)
}
