package api

import (
	"encoding/json"
	"fmt"
	"github.com/panjiawan/go-lib/pkg/plog"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"websocket/conf"
	"websocket/protoc/pb"
	"websocket/service/code"
	"websocket/service/dto"
	"websocket/service/internal"
)

func GrpcSend(ctx *fasthttp.RequestCtx) {
	req := &dto.WebSocketSendReq{}
	body := ctx.PostBody()
	if err := json.Unmarshal(body, &req); err != nil {
		plog.Error("GrpcSend json.Unmarshal error", zap.Error(err), zap.Any("body", string(body)))
		internal.OutputError(ctx, code.ErrorParameter)
		return
	}
	grpcAddr := fmt.Sprintf("127.0.0.1:%d", conf.GetHandle().GetHttpConf().GrpcPort)
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		plog.Error("grpc dial failure", zap.String("addr", grpcAddr), zap.Error(err))
		return
	}
	defer conn.Close()

	client := pb.NewWebsocketClient(conn)
	_, err = client.WebSocketSend(ctx, &pb.WebSocketSendReq{
		Platform: req.Platform,
		ToUsers:  req.ToUsers,
		Message:  req.Message,
	})

	if err != nil {
		plog.Error("GrpcSend error", zap.Error(err), zap.Any("req", req))
		internal.OutputError(ctx, code.ErrorServer)
		return
	}

	internal.OutputSuccess(ctx)
}
