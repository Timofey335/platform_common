package kafka

import (
	"context"

	"github.com/Timofey335/auth/internal/client/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
