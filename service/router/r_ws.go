package router

import (
	"github.com/panjiawan/go-lib/pkg/phttp"
	"websocket/service/control/api"
)

func init() {
	groupRoute := NewGroup("api/ws")
	groupRoute.Add("sub", &routerMethod{api.WsSubEvent, phttp.MethodGet, false})      //连接  // nchan 的订阅事件回调
	groupRoute.Add("unsub", &routerMethod{api.WsUnSubEvent, phttp.MethodPost, false}) //断开连接  // nchan 的取消订阅事件回调
}
