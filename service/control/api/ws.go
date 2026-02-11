package api

import (
	"encoding/json"
	"github.com/panjiawan/go-lib/pkg/plog"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"websocket/service/code"
	"websocket/service/control"
	"websocket/service/defines"
	"websocket/service/dto"
	"websocket/service/internal"
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

func IsOnline(ctx *fasthttp.RequestCtx) {
	req := &dto.IsOnlineUserReq{}
	body := ctx.PostBody()
	if err := json.Unmarshal(body, &req); err != nil {
		plog.Error("IsOnline json.Unmarshal error", zap.Error(err), zap.Any("body", string(body)))
		internal.OutputError(ctx, code.ErrorParameter)
		return
	}
	// 验证数据
	if req.Platform == "" || len(req.UserIds) <= 0 {
		plog.Error("IsOnline Check req empty", zap.Any("body", string(body)))
		internal.OutputError(ctx, code.ErrorParameter)
		return
	}

	users, err := control.IsOnline(ctx, req)
	if err != nil {
		plog.Error("IsOnline error", zap.Error(err), zap.Any("req", req))
		internal.OutputError(ctx, code.ErrorServer)
		return
	}

	internal.Output(ctx, users)
	return
}

func GetOnlineUserCount(ctx *fasthttp.RequestCtx) {
	req := &dto.GetOnlineUserCountReq{}
	body := ctx.PostBody()
	if err := json.Unmarshal(body, &req); err != nil {
		plog.Error("GetOnlineUserCount json.Unmarshal error", zap.Error(err), zap.Any("body", string(body)))
		internal.OutputError(ctx, code.ErrorParameter)
		return
	}
	// 验证数据
	if req.Platform == "" {
		plog.Error("GetOnlineUserCount Check req empty", zap.Any("body", string(body)))
		internal.OutputError(ctx, code.ErrorParameter)
		return
	}
	count, err := control.GetOnlineUserCount(ctx, req)
	if err != nil {
		plog.Error("GetOnlineUserCount error", zap.Error(err), zap.Any("req", req))
		internal.OutputError(ctx, code.ErrorServer)
		return
	}

	internal.Output(ctx, count)
	return
}

func WebSocketSend(ctx *fasthttp.RequestCtx) {
	req := &dto.WebSocketSendReq{}
	body := ctx.PostBody()
	if err := json.Unmarshal(body, &req); err != nil {
		plog.Error("WebSocketSend json.Unmarshal error", zap.Error(err), zap.Any("body", string(body)))
		internal.OutputError(ctx, code.ErrorParameter)
		return
	}
	err := control.SendMessage(ctx, req)
	if err != nil {
		plog.Error("WebSocketSend error", zap.Error(err), zap.Any("req", req))
		internal.OutputError(ctx, code.ErrorServer)
		return
	}
	plog.Debug("WebSocketSend", zap.Any("req", req))

	internal.OutputSuccess(ctx)
}

func GetPing(ctx *fasthttp.RequestCtx) {
	err := control.GetPing(ctx)
	plog.Debug("GetPing。。")
	if err != nil {
		plog.Error("GetPing error", zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		internal.OutputError(ctx, code.ErrorServer)
		return
	}

	internal.OutputSuccess(ctx)
}

func WSSocketSend(ctx *fasthttp.RequestCtx) {
	req := &dto.WebSocketSendReq{}
	body := ctx.PostBody()
	plog.Debug("WSSocketSend", zap.Any("body", string(body)))
	if err := json.Unmarshal(body, &req); err != nil {
		plog.Error("WebSocketSend json.Unmarshal error", zap.Error(err), zap.Any("body", string(body)))
		internal.OutputError(ctx, code.ErrorParameter)
		return
	}
	err := control.SendMessage(ctx, req)
	if err != nil {
		plog.Error("WebSocketSend error", zap.Error(err), zap.Any("req", req))
		internal.OutputError(ctx, code.ErrorServer)
		return
	}

	internal.OutputSuccess(ctx)
}
