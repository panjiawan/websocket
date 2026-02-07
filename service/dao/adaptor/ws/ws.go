package ws

import (
	"context"
	"websocket/service/dto"
)

type IWebsocket interface {
	SendWebsocket(ctx context.Context, msgReq *dto.WebSocketSendReq) error
}
