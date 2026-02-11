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

func (s *WsGrpcService) IsOnlineUsers(ctx context.Context, req *pb.IsOnlineUsersReq) (*pb.IsOnlineUsersResp, error) {
	uIDs, err := control.IsOnline(ctx, &dto.IsOnlineUserReq{
		Platform: req.Platform,
		UserIds:  req.UserIds,
	})
	if err != nil {
		plog.Error("IsOnlineUsers error", zap.Error(err), zap.Any("req", req))
		return nil, errors.New(code.ErrorServer.Msg)
	}

	return &pb.IsOnlineUsersResp{
		Users: uIDs,
	}, nil
}

func (s *WsGrpcService) GetOnlineCount(ctx context.Context, req *pb.GetOnlineCountReq) (*pb.GetOnlineCountResp, error) {
	count, err := control.GetOnlineUserCount(ctx, &dto.GetOnlineUserCountReq{
		Platform: req.Platform,
	})
	if err != nil {
		plog.Error("GetOnlineCount error", zap.Error(err), zap.Any("req", req))
		return nil, errors.New(code.ErrorServer.Msg)
	}

	return &pb.GetOnlineCountResp{
		Count: count,
	}, nil
}
