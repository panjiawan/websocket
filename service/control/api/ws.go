package api

import (
	"github.com/panjiawan/go-lib/pkg/plog"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"websocket/service/control"
	"websocket/service/defines"
	"websocket/service/dto"
)

func WsSubEvent(ctx *fasthttp.RequestCtx) {
	channel := ctx.Request.Header.Peek("X-Channel-Id")
	platform := ctx.Request.Header.Peek("X-Channel-Group")
	plog.Debug("WsSubEvent", zap.String("channel", string(channel)), zap.String("platform", string(platform)))

	control.NChanEvent(ctx, &dto.NChanEventDto{
		EventType: defines.WsSubscribeEvent,
		ChannelID: string(channel),
		Platform:  string(platform),
	})

	//ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.WriteString("ok")
}

func WsUnSubEvent(ctx *fasthttp.RequestCtx) {
	channel := ctx.Request.Header.Peek("X-Channel-Id")
	platform := ctx.Request.Header.Peek("X-Channel-Group")
	plog.Debug("WsUnSubEvent", zap.String("channel", string(channel)), zap.String("platform", string(platform)))

	control.NChanEvent(ctx, &dto.NChanEventDto{
		EventType: defines.WsUnSubscribeEvent,
		ChannelID: string(channel),
		Platform:  string(platform),
	})
	ctx.WriteString("ok")
}
