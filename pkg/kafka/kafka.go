package kafka

import (
	"context"

	"github.com/Timofey335/platform_common/pkg/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
