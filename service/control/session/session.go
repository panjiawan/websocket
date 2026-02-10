package session

import (
	"context"
	"github.com/panjiawan/go-lib/pkg/plog"
	"go.uber.org/zap"
	"sync"
	"websocket/service/dao/adaptor/redis"
	"websocket/service/defines"
)

type Session struct {
	sessionDrive redis.ISession
	lock         sync.RWMutex
	onlineMap    map[string]map[string]bool // map[platform]map[user_id]bool
}

func NewSession() *Session {
	onlineMap := make(map[string]map[string]bool)
	for _, v := range defines.PlatformAll {
		onlineMap[v] = make(map[string]bool)
	}
	return &Session{
		onlineMap:    onlineMap,
		sessionDrive: redis.NewSession(),
	}
}

//订阅
func (s *Session) Subscribe(ctx context.Context, platform, userID string) error {
	plog.Debug("Subscribe", zap.String("userID", userID), zap.String("platform", platform))

	err := s.sessionDrive.AddOnline(ctx, platform, userID)
	if err != nil {
		plog.Error("Subscribe error", zap.String("userID", userID), zap.String("platform", platform), zap.Error(err))
		return err
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.onlineMap[platform][userID] = true
	return nil
}

//取消订阅
func (s *Session) UnSubscribe(ctx context.Context, platform, userID string) error {
	plog.Debug("UnSubscribe", zap.String("userID", userID), zap.String("platform", platform))
	err := s.sessionDrive.RemOnline(ctx, platform, userID)
	if err != nil {
		plog.Error("UnSubscribe error", zap.String("userID", userID), zap.String("platform", platform), zap.Error(err))
		return err
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.onlineMap[platform]; !ok {
		delete(s.onlineMap[platform], userID)
	}
	return nil
}

//节点心跳
func (s *Session) NodeHeartbeat(ctx context.Context) error {
	return s.sessionDrive.NodeHeartbeat(ctx)
}

func (s *Session) CleanOfflineNodes(ctx context.Context) error {
	return s.sessionDrive.CleanOfflineNode(ctx)
}

func (s *Session) CleanSelfSession(ctx context.Context) {
	s.sessionDrive.CleanSelfNode(ctx)

	s.lock.Lock()
	defer s.lock.Unlock()
	s.onlineMap = make(map[string]map[string]bool)
}

func (s *Session) CheckInSession(ctx context.Context, platform string, userIds []string) map[string]bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	retMap := make(map[string]bool)
	platformMap := s.onlineMap[platform]
	for _, v := range userIds {
		retMap[v] = platformMap[v]
	}
	return retMap
}

func (s *Session) GetNodeOnlineUserIds(ctx context.Context, platform string) []string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	userIds := make([]string, 0)
	platformMap := s.onlineMap[platform]
	for k, v := range platformMap {
		if v {
			userIds = append(userIds, k)
		}
	}
	return userIds
}

func (s *Session) IsOnline(ctx context.Context, platform string, userIDs []string) (map[string]bool, error) {
	return s.sessionDrive.IsOnline(ctx, platform, userIDs)
}

func (s *Session) GetOnlineCount(ctx context.Context, platform string) (int64, error) {
	return s.sessionDrive.GetOnlineCount(ctx, platform)
}

func (s *Session) GetPing(ctx context.Context) error {
	return s.sessionDrive.GetPing(ctx)
}
