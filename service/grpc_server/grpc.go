package grpc_server

import (
	"fmt"
	"github.com/panjiawan/go-lib/pkg/plog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"runtime/debug"
	"websocket/conf"
	"websocket/protoc/pb"
)

type WsGrpcService struct {
	pb.UnimplementedWebsocketServer
}

func Run() {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				plog.Error("grpc server error", zap.Error(fmt.Errorf("%v", e)), zap.String("panic", string(debug.Stack())))
			}
		}()
		addr := fmt.Sprintf(":%d", conf.GetHandle().GetHttpConf().GrpcPort)
		srv := grpc.NewServer()
		// 注册服务
		pb.RegisterWebsocketServer(srv, &WsGrpcService{})
		// 启用 gRPC 反射
		reflection.Register(srv)

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			plog.Error("grpc start error", zap.Error(err))
			panic("grpc start error")
		}
		plog.Info("grpc server start", zap.String("addr", addr))
		err = srv.Serve(listener)

		if err != nil {
			plog.Error("grpc start error", zap.Error(err))
			panic("grpc start error")
		}
	}()
}
