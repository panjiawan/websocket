package router

import (
	"websocket/service/control/api"
)

var postHandleList = map[string]*routerMethod{
	"/": {
		Handle: api.Index,
		Filter: false,
	},
}
