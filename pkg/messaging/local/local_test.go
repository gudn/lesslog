package local_test

import (
	"context"
	"sync"
	"testing"

	"github.com/gudn/lesslog/pkg/messaging/local"
	"github.com/stretchr/testify/assert"
)

func collectChan(t *testing.T, ctx context.Context, ch <-chan uint64) []uint64 {
	if deadline, ok := t.Deadline(); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}
	result := make([]uint64, 0)
	for {
		select {
		case <-ctx.Done():
			return result
		case val, ok := <-ch:
			if !ok {
				return result
			}
			result = append(result, val)
		}
	}
}

func testChan(
	t *testing.T,
	wg *sync.WaitGroup,
	ctx context.Context,
	ch <-chan uint64,
	expected []uint64,
	msg ...any,
) {
	defer wg.Done()
	got := collectChan(t, ctx, ch)
	assert.Equal(t, got, expected, msg...)
}

func TestLocal(t *testing.T) {
	ctx := context.Background()
	nested, cancel := context.WithCancel(ctx)
	m := local.New()
	ch11, _ := m.Listen(nested, "log1")
	ch12, _ := m.Listen(nested, "log1")
	ch2, _ := m.Listen(nested, "log2")
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go testChan(t, wg, ctx, ch11, []uint64{1, 2, 1}, "log1 (1)")
	go testChan(t, wg, ctx, ch12, []uint64{1, 2, 1}, "log1 (2)")
	go testChan(t, wg, ctx, ch2, []uint64{3, 3}, "log2 (1)")
	m.Post(ctx, "log1", 1)
	m.Post(ctx, "log2", 3)
	m.Post(ctx, "log1", 2)
	m.Post(ctx, "log1", 1)
	m.Post(ctx, "log2", 3)
	cancel()
	wg.Wait()
}
