package local

import (
	"context"
	"sync"
)

type LocalMessaging struct {
	sync.RWMutex
	chans map[string]map[*chan uint64]struct{}
}

func (m *LocalMessaging) Listen(ctx context.Context, log_name string) (
	<-chan uint64,
	error,
) {
	m.Lock()
	defer m.Unlock()
	ch := make(chan uint64)

	if val, ok := m.chans[log_name]; ok {
		val[&ch] = struct{}{}
	} else {
		val := make(map[*chan uint64]struct{})
		val[&ch] = struct{}{}
		m.chans[log_name] = val
	}

	go func() {
		select {
		case <-ctx.Done():
			m.Lock()
			defer m.Unlock()
			if val, ok := m.chans[log_name]; ok {
				delete(val, &ch)
				close(ch)
				return
			}
		}
	}()

	return ch, nil
}

func (m *LocalMessaging) Post(
	_ context.Context,
	log_name string,
	value uint64,
) error {
	m.RLock()
	defer m.RUnlock()
	if val, ok := m.chans[log_name]; ok {
		for ch := range val {
			*ch <- value
		}
	}
	return nil
}

func New() *LocalMessaging {
	return &LocalMessaging{chans: make(map[string]map[*chan uint64]struct{})}
}
