package control

import (
	"context"
	"github.com/panjiawan/go-lib/pkg/plog"
	"go.uber.org/zap"
	"websocket/service/dto"
)

// 2. 获取用户是否在线
func IsOnline(ctx context.Context, req *dto.IsOnlineUserReq) ([]string, error) {
	onlineMap, err := servicer.session.IsOnline(ctx, req.Platform, req.UserIds)
	if err != nil {
		plog.Error("IsOnline IsOnline error", zap.Error(err), zap.Any("req", req))
		return nil, err
	}
	userIds := make([]string, 0)
	for uid, v := range onlineMap {
		if v {
			userIds = append(userIds, uid)
		}
	}
	return userIds, nil
}

// 3. 获取在线用户数
func GetOnlineUserCount(ctx context.Context, req *dto.GetOnlineUserCountReq) (int64, error) {
	count, err := servicer.session.GetOnlineCount(ctx, req.Platform)
	if err != nil {
		plog.Error("GetOnlineUserCount GetOnlineCount error", zap.Error(err), zap.Any("req", req))
		return 0, err
	}
	return count, nil
}
