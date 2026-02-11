package router

import (
	"github.com/panjiawan/go-lib/pkg/phttp"
	"websocket/service/control/api"
)

func init() {
	groupRoute := NewGroup("api/ws")
	groupRoute.Add("sub", &routerMethod{api.WsSubEvent, phttp.MethodGet, false})     //连接  // nchan 的订阅事件回调
	groupRoute.Add("unsub", &routerMethod{api.WsUnSubEvent, phttp.MethodGet, false}) //断开连接  // nchan 的取消订阅事件回调

	groupRoute.Add("online/check", &routerMethod{api.IsOnline, phttp.MethodPost, false})           // 检查用户是否在线
	groupRoute.Add("online/count", &routerMethod{api.GetOnlineUserCount, phttp.MethodPost, false}) // 获取在线用户数

	groupRoute.Add("send", &routerMethod{api.WebSocketSend, phttp.MethodPost, false}) // 发布websocket消息

	groupRoute.Add("ws_send", &routerMethod{api.WSSocketSend, phttp.MethodPost, false}) // WS连接发送的消息

	groupRoute.Add("grpc_send", &routerMethod{api.GrpcSend, phttp.MethodPost, false}) // 发布websocket消息
}
