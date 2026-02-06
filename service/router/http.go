package router

import (
	"fmt"
	"github.com/panjiawan/go-lib/pkg/phttp"
	"github.com/panjiawan/go-lib/pkg/plog"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"runtime/debug"
	"websocket/conf"
	"websocket/service/code"
	"websocket/service/internal"
	"websocket/service/middleware"
)

type HttpRouter struct {
	httpConf *conf.HttpConf
	handle   *phttp.Service
}

type routerMethod struct {
	Handle func(ctx *fasthttp.RequestCtx)
	Method string
	Filter bool
}

func New(cfg *conf.HttpConf) *HttpRouter {
	return &HttpRouter{
		httpConf: cfg,
	}
}

// Run 启动函数
func (h *HttpRouter) Run() {
	h.handle = phttp.New(
		phttp.WithAddress("", h.httpConf.HttpPort),
		phttp.WithCertificate(h.httpConf.HttpsCertFile, h.httpConf.HttpsKeyFile),
		phttp.WithRate(h.httpConf.RateLimitPerSec, h.httpConf.RateLimitCapacity),
	)

	h.Register()

	if err := h.handle.Run(); err != nil {
		plog.Error("http server start error", zap.Error(err))
	}
}

func (h *HttpRouter) Register() {
	parseGroup()
	for path, v := range routesList {
		h.handle.Register(path, v.Method, h.PrepareCall)
		plog.Info("register", zap.String("path", path), zap.String("method", v.Method))
	}
}

func (h *HttpRouter) PrepareCall(ctx *fasthttp.RequestCtx) {
	defer func() {
		if e := recover(); e != nil {
			plog.Error("panic prepareCall", zap.Error(fmt.Errorf("%v", e)), zap.String("trace", string(debug.Stack())))
		}
	}()

	h.options(ctx)
	path := string(ctx.URI().Path())
	//method := string(ctx.Request.Header.Method())

	if _, ok := routesList[path]; ok {
		if routesList[path].Filter {
			if outputCode := h.Filter(ctx); outputCode != code.Success {
				internal.OutputError(ctx, outputCode)
				return
			}
		}
		routesList[path].Handle(ctx)
	}
}

func (h *HttpRouter) options(ctx *fasthttp.RequestCtx) {
	// 处理OPTIONS
	middleware.SetCORSHeader(ctx)
	ctx.SetStatusCode(fasthttp.StatusAccepted)
}

func (h *HttpRouter) Close() {
}
