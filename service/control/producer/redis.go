package producer

import (
	"context"
	"websocket/service/dao/adaptor/redis"
	"websocket/service/dto"
)

type RedisProducer struct {
	publisher redis.IMessage
}

func NewRedisProducer() *RedisProducer {
	return &RedisProducer{
		publisher: redis.NewMessage(),
	}
}

func (p *RedisProducer) PublishMessage(ctx context.Context, req *dto.WebSocketSendReq) error {
	return p.publisher.PublishMessage(ctx, req)
}
