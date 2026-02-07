package consumer

import (
	"context"
	"encoding/json"
	"github.com/panjiawan/go-lib/pkg/plog"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"sync"
	"time"
	"websocket/conf"
	"websocket/service/control/session"
	"websocket/service/dao/adaptor/redis"
	"websocket/service/dao/adaptor/ws"
	"websocket/service/dto"
)

type RedisConsumer struct {
	session    *session.Session
	subscriber redis.IMessage
	wsSender   ws.IWebsocket
}

func NewRedisConsumer(session *session.Session) *RedisConsumer {
	return &RedisConsumer{
		session:    session,
		subscriber: redis.NewMessage(),
		wsSender:   ws.NewNChanSend(),
	}
}

func (c *RedisConsumer) ConsumeMessage(ctx context.Context, platform string) error {
	var (
		wg      sync.WaitGroup
		msgChan = make(chan *dto.WebSocketSendReq)
	)
	c.sendWebsocket(ctx, msgChan, &wg)
	defer wg.Wait()

	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					plog.Error("ConsumeMessage panic recovered ",
						zap.Any("err", err), zap.Any("platform", platform))
				}
			}()
			sub := c.subscriber.SubscribeMessage(ctx, platform)
			for {
				msg, err := sub.ReceiveMessage(ctx)
				if err != nil {
					plog.Error("ConsumeMessage receive message error",
						zap.Error(err), zap.Any("platform", platform))
					break
				}
				payload := msg.Payload

				wsMsgReq := &dto.WebSocketSendReq{}
				err = json.Unmarshal([]byte(payload), wsMsgReq)
				if err != nil {
					plog.Error("ConsumeMessage json unmarshal error",
						zap.Error(err), zap.Any("payload", payload))
					continue
				}
				if len(wsMsgReq.ToUsers) == 0 {
					// 说明是广播，也就是相当于把自己的所有session中的对象都取出来，并进行发送
					userIds := c.session.GetNodeOnlineUserIds(ctx, wsMsgReq.Platform)
					wsMsgReq.ToUsers = userIds
				} else {
					// 说明指定user_id发送，但是有些user_id他不在这个节点，要进行过滤
					onlineMap := c.session.CheckInSession(ctx, wsMsgReq.Platform, wsMsgReq.ToUsers)
					wsMsgReq.ToUsers = lo.Filter(wsMsgReq.ToUsers, func(item string, index int) bool {
						return onlineMap[item]
					})
				}
				if len(wsMsgReq.ToUsers) == 0 {
					continue
				}
				msgChan <- wsMsgReq
			}
		}()

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second * 1):
			continue
		}
	}
}

func (c *RedisConsumer) sendWebsocket(ctx context.Context, msgChan chan *dto.WebSocketSendReq, wg *sync.WaitGroup) {
	for i := 0; i < conf.GetHandle().GetHttpConf().MaxGoCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-msgChan:
					err := c.wsSender.SendWebsocket(ctx, msg)
					if err != nil {
						plog.Error("send websocket error", zap.Error(err), zap.Any("msg", msg))
						continue
					}
				}
			}
		}()
	}
}
