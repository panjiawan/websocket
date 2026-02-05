package internal

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"time"
	"websocket/service/code"
	"websocket/service/middleware"
)

type ResData map[string]interface{}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Time int64       `json:"time,omitempty"`
}

func Output(ctx *fasthttp.RequestCtx, data interface{}) {
	middleware.OutputMiddleware(ctx)

	res := &Response{
		Code: code.Success.Code,
		Data: data,
		Time: time.Now().UnixMilli(),
	}
	resByte, err := jsoniter.Marshal(res)

	if err != nil {
		return
	}

	ctx.Write(resByte)
}

func OutputError(ctx *fasthttp.RequestCtx, code code.OutputCode) {
	middleware.OutputMiddleware(ctx)
	res := &Response{
		Code: code.Code,
		Msg:  code.Msg,
		Time: time.Now().UnixMilli(),
	}
	resByte, err := jsoniter.Marshal(res)

	if err != nil {
		return
	}

	ctx.Write(resByte)
}

func OutputCustomError(ctx *fasthttp.RequestCtx, filedDesc string) {
	middleware.OutputMiddleware(ctx)
	res := &Response{
		Code: code.ErrorParam.Code,
		Msg:  filedDesc,
		Time: time.Now().UnixMilli(),
	}
	resByte, err := jsoniter.Marshal(res)

	if err != nil {
		return
	}

	ctx.Write(resByte)
}

func OutputSuccess(ctx *fasthttp.RequestCtx) {
	middleware.OutputMiddleware(ctx)

	res := &Response{
		Code: code.Success.Code,
		Time: time.Now().UnixMilli(),
	}
	resByte, err := jsoniter.Marshal(res)

	if err != nil {
		return
	}

	ctx.Write(resByte)
}
