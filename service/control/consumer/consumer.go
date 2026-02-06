package consumer

import "context"

//消费者
type IConsumer interface {
	ConsumeMessage(ctx context.Context, platform string) error
}
