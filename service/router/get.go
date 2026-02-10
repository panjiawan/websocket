package router

import (
	"websocket/service/control/api"
)

var getHandleList = map[string]*routerMethod{
	"/": {
		Handle: api.Index,
		Filter: false,
	},
	"/ping": {
		Handle: api.GetPing,
		Filter: false,
	},
}
