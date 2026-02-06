package defines

import "time"

// 增加/减少
const (
	Increment = 1
	Decrement = -1
)

const (
	WsSubscribeEvent   = "subscribe"
	WsUnSubscribeEvent = "unsubscribe"
)

const (
	PlatformWeb      = "web"
	PlatformApp      = "app"
	PlatformAdmin    = "admin"
	PlatformCustomer = "customer"
)

var PlatformAll = []string{
	PlatformWeb,
	PlatformApp,
	PlatformAdmin,
	PlatformCustomer,
}

const (
	NodeHeartbeatTimeout = time.Minute * 5
)
