package router

import (
	"github.com/valyala/fasthttp"
	"websocket/service/code"
)

func (h *HttpRouter) Filter(ctx *fasthttp.RequestCtx) code.OutputCode {
	//if internal.VerifyAuth(ctx) {
	//	return code.Success
	//}

	//return code.ErrorAuth

	return code.Success
}
