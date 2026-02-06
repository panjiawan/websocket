package producer

import (
	"context"
	"websocket/service/dto"
)

//生产者
type IProducer interface {
	PublishMessage(ctx context.Context, req *dto.WebSocketSendReq) error
}
