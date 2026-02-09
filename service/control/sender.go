package control

import (
	"context"
	"github.com/panjiawan/go-lib/pkg/plog"
	"go.uber.org/zap"
	"websocket/service/dto"
)

// 发送websocket消息
func SendMessage(ctx context.Context, req *dto.WebSocketSendReq) error {
	err := servicer.producer.PublishMessage(ctx, req)
	if err != nil {
		plog.Error("SendMessage PublishMessage error", zap.Error(err), zap.Any("req", req))
		return err
	}
	return nil
}
