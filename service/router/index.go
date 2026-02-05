package router

import (
	"github.com/panjiawan/go-lib/pkg/phttp"
	"websocket/service/control/api"
)

func init() {
	groupRoute := NewGroup("index")
	groupRoute.Add("get", &routerMethod{api.Index, phttp.MethodGet, false})
	groupRoute.Add("post", &routerMethod{api.Index, phttp.MethodPost, false})
}
