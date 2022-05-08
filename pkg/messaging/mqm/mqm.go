package mqm

import (
	"context"
	"encoding/binary"
	"errors"

	"github.com/gudn/lesslog/internal/mq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var connIsNil error = status.Error(codes.Unavailable, "nats is unavailable")

type MqMessaging struct{}

func (MqMessaging) Listen(
	ctx context.Context,
	log_name string,
) (<-chan uint64, error) {
	if mq.Conn == nil {
		return nil, connIsNil
	}

	sub, err := mq.Conn.SubscribeSync(log_name)
	if err != nil {
		return nil, err
	}

	result := make(chan uint64)
	go func() {
		defer close(result)
		for {
			msg, err := sub.NextMsgWithContext(ctx)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Error().Err(err).Msg("failed receive message")
				}
				return
			}
			val, n := binary.Uvarint(msg.Data)
			if n > 0 {
				result <- val
			} else {
				log.Warn().Int("readed", n).Msg("failed to parse message")
			}
		}
	}()

	return result, nil
}

func (MqMessaging) Post(
	_ context.Context,
	log_name string,
	value uint64,
) error {
	if mq.Conn == nil {
		return connIsNil
	}
	buf := make([]byte, 10)
	used := binary.PutUvarint(buf, value)
	return mq.Conn.Publish(log_name, buf[:used])
}
