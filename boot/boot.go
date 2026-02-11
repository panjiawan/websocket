package boot

import (
	"github.com/panjiawan/go-lib/pkg/plog"
	"websocket/conf"
	"websocket/service/control"
	"websocket/service/dao"
	"websocket/service/grpc_server"
	"websocket/service/router"
)

type BootArgs struct {
	EtcPath string
	LogPath string
}

type BootHandler interface {
	Run()
	Close()
}

func Start(etcPath string, logPath string) {
	// load conf
	confHandle := conf.New(etcPath)
	confHandle.Run()
	// start log
	plog.Start(logPath, "websocket_log", confHandle.GetHttpConf().EnableDebug, confHandle.GetHttpConf().EnableStdout)

	// start signal
	go closeSignalListen()

	plog.Info("conf started")
	dao.Run()
	plog.Info("model started")
	/////////////////////
	control.Run()
	defer control.Stop()
	////////////////////

	//启动grpc服务
	grpc_server.Run()

	route := router.New(confHandle.GetHttpConf())
	route.Run()
}

func Stay() {
	select {}
}

// 优雅关闭调用点
func close() {
}
