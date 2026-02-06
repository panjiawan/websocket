package api

import (
	"github.com/panjiawan/go-lib/pkg/plog"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func WsSubEvent(ctx *fasthttp.RequestCtx) {
	channel := ctx.Request.Header.Peek("X-Channel-Id")
	platform := ctx.Request.Header.Peek("X-Channel-Group")

	plog.Info("WsSubEvent", zap.String("channel", string(channel)), zap.String("platform", string(platform)))
}

func WsUnSubEvent(ctx *fasthttp.RequestCtx) {
	channel := ctx.Request.Header.Peek("X-Channel-Id")
	platform := ctx.Request.Header.Peek("X-Channel-Group")

	plog.Info("WsUnSubEvent", zap.String("channel", string(channel)), zap.String("platform", string(platform)))
}
