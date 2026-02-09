package control

import (
	"context"
	"github.com/panjiawan/go-lib/pkg/plog"
	"go.uber.org/zap"
	"websocket/service/defines"
	"websocket/service/dto"
)

func NChanEvent(ctx context.Context, req *dto.NChanEventDto) {
	var err error
	switch req.EventType {
	case defines.WsSubscribeEvent:
		err = servicer.session.Subscribe(ctx, req.Platform, req.ChannelID)
	case defines.WsUnSubscribeEvent:
		err = servicer.session.UnSubscribe(ctx, req.Platform, req.ChannelID)
	default:
		plog.Info("NChanEvent unknown event type", zap.Any("req", req))
	}
	if err != nil {
		plog.Error("NChanEvent error", zap.Error(err), zap.Any("req", req))
	}
}
