package redis

import (
	"context"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/redis/go-redis/v9"
	"websocket/service/dao"
	"websocket/service/dto"
)

type IMessage interface {
	PublishMessage(ctx context.Context, req *dto.WebSocketSendReq) error
	SubscribeMessage(ctx context.Context, platform string) *redis.PubSub
}

type Message struct {
	client *redis.Client
}

func NewMessage() *Message {
	return &Message{
		client: dao.Redis(dao.RedisKey),
	}
}

func fmtSubscribeRedisKey(platform string) string {
	return dao.FormatRedisKey(fmt.Sprintf("ws:message:%s", platform))
}

func (m *Message) PublishMessage(ctx context.Context, req *dto.WebSocketSendReq) error {
	redisKey := fmtSubscribeRedisKey(req.Platform)
	// TODO tips 确认下内容是否最后是json字符串
	return m.client.Publish(ctx, redisKey, gconv.String(req)).Err()
}

func (m *Message) SubscribeMessage(ctx context.Context, platform string) *redis.PubSub {
	redisKey := fmtSubscribeRedisKey(platform)
	return m.client.Subscribe(ctx, redisKey)
}
