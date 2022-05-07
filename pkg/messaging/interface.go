package messaging

import "context"

type Interface interface {
	Listen(ctx context.Context, log_name string) (<-chan uint64, error)
	Post(ctx context.Context, log_name string, value uint64) error
}
