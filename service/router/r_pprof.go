package router

import (
	"github.com/panjiawan/go-lib/pkg/phttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"net/http/pprof"
)

func init() {
	groupRoute := NewGroup("debug/pprof")

	groupRoute.Add("index", &routerMethod{fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index), phttp.MethodGet, false})
	groupRoute.Add("cmdline", &routerMethod{fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Cmdline), phttp.MethodGet, false})
	groupRoute.Add("profile", &routerMethod{fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile), phttp.MethodGet, false})
	groupRoute.Add("symbol", &routerMethod{fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Symbol), phttp.MethodGet, false})
	groupRoute.Add("trace", &routerMethod{fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Trace), phttp.MethodGet, false})

	groupRoute.Add("heap", &routerMethod{fasthttpadaptor.NewFastHTTPHandler(pprof.Handler("heap")), phttp.MethodGet, false})
	groupRoute.Add("goroutine", &routerMethod{fasthttpadaptor.NewFastHTTPHandler(pprof.Handler("goroutine")), phttp.MethodGet, false})
	groupRoute.Add("block", &routerMethod{fasthttpadaptor.NewFastHTTPHandler(pprof.Handler("block")), phttp.MethodGet, false})
	groupRoute.Add("mutex", &routerMethod{fasthttpadaptor.NewFastHTTPHandler(pprof.Handler("mutex")), phttp.MethodGet, false})

}
