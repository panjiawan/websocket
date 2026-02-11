package grpc_server

import (
	"context"
	"errors"
	"github.com/panjiawan/go-lib/pkg/plog"
	"go.uber.org/zap"
	"websocket/protoc/pb"
	"websocket/service/code"
	"websocket/service/control"
	"websocket/service/dto"
)

// 发送websocket消息

func (s *WsGrpcService) WebSocketSend(ctx context.Context, req *pb.WebSocketSendReq) (*pb.WebSocketSendResp, error) {
	err := control.SendMessage(ctx, &dto.WebSocketSendReq{
		Message:  req.Message,
		Platform: req.Platform,
		ToUsers:  req.ToUsers,
	})
	if err != nil {
		plog.Error("WebSocketSend error", zap.Error(err), zap.Any("req", req))
		return nil, errors.New(code.ErrorServer.Msg)
	}

	return &pb.WebSocketSendResp{}, nil
}
