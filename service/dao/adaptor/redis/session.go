package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"time"
	"websocket/service/dao"
	"websocket/service/defines"
	"websocket/service/utils"
)

type ISession interface {
	AddOnline(ctx context.Context, platform string, userId string) error
	RemOnline(ctx context.Context, platform string, userId string) error
	IsOnline(ctx context.Context, platform string, userIds []string) (map[string]bool, error)
	GetOnlineCount(ctx context.Context, platform string) (int64, error)
	//TODO 在线用户列表

	NodeHeartbeat(ctx context.Context) error
	CleanSelfNode(ctx context.Context) error
	CleanOfflineNode(ctx context.Context) error
	GetPing(ctx context.Context) error
}

type Session struct {
	nodeID string
	client *redis.Client
}

func NewSession() *Session {
	return &Session{
		nodeID: utils.UUIDHex(),
		client: dao.Redis(dao.RedisKey),
	}
}

func fmtOnlineRedisKey(platform string, nodeID string) string {
	return dao.FormatRedisKey(fmt.Sprintf("ws:%s:%s:online:set",
		platform, nodeID))
}

func fmtNodeHeartbeatRedisKey() string {
	return dao.FormatRedisKey("ws:node:heartbeat:hash")
}

func (s *Session) AddOnline(ctx context.Context, platform string, userId string) error {
	redisKey := fmtOnlineRedisKey(platform, s.nodeID)
	_, err := s.client.SAdd(ctx, redisKey, userId).Result()
	return err
}
func (s *Session) RemOnline(ctx context.Context, platform string, userId string) error {
	redisKey := fmtOnlineRedisKey(platform, s.nodeID)
	_, err := s.client.SRem(ctx, redisKey, userId).Result()
	return err
}

func (s *Session) IsOnline(ctx context.Context, platform string, userIds []string) (map[string]bool, error) {
	userIds = lo.Uniq(userIds) //去重
	if len(userIds) > 500 {
		return nil, errors.New("用户ID数量超出限制,max=500")
	}
	retMap := make(map[string]bool)
	hashKey := fmtNodeHeartbeatRedisKey()
	nodeHeartbeatMap, err := s.client.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return nil, err
	}
	for nodeID, _ := range nodeHeartbeatMap {
		rdsKey := fmtOnlineRedisKey(platform, nodeID)
		pipe := s.client.Pipeline()
		cmdS := make([]*redis.BoolCmd, len(userIds))
		for i, m := range userIds {
			cmdS[i] = pipe.SIsMember(ctx, rdsKey, m)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			return nil, err
		}
		for i, uid := range userIds {
			exists := cmdS[i].Val()
			if exists {
				retMap[uid] = true
			}
		}
	}
	return retMap, nil
}

func (s *Session) GetOnlineCount(ctx context.Context, platform string) (int64, error) {
	hashKey := fmtNodeHeartbeatRedisKey()
	nodeHeartbeatMap, err := s.client.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return 0, err
	}
	count := int64(0)
	for nodeID, _ := range nodeHeartbeatMap {
		rdsKey := fmtOnlineRedisKey(platform, nodeID)
		temp, err := s.client.SCard(ctx, rdsKey).Result()
		if err != nil {
			return 0, err
		}
		count += temp
	}
	return count, nil
}

func (s *Session) NodeHeartbeat(ctx context.Context) error {
	hashKey := fmtNodeHeartbeatRedisKey()
	_, err := s.client.HSet(ctx, hashKey, s.nodeID, time.Now().Unix()).Result()
	return err
}
func (s *Session) CleanSelfNode(ctx context.Context) error {
	hashKey := fmtNodeHeartbeatRedisKey()
	pipe := s.client.Pipeline()
	pipe.HDel(ctx, hashKey, s.nodeID)
	for _, platform := range defines.PlatformAll {
		pipe.Del(ctx, fmtOnlineRedisKey(platform, s.nodeID))
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (s *Session) CleanOfflineNode(ctx context.Context) error {
	hashKey := fmtNodeHeartbeatRedisKey()
	nodeHeartbeatMap, err := s.client.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return err
	}
	for nodeID, heartBeatStr := range nodeHeartbeatMap {
		heartBeatTime := time.Unix(gconv.Int64(heartBeatStr), 0)
		if time.Now().Sub(heartBeatTime) > defines.NodeHeartbeatTimeout {
			pipe := s.client.Pipeline()
			pipe.HDel(ctx, hashKey, nodeID)
			for _, platform := range defines.PlatformAll {
				pipe.Del(ctx, fmtOnlineRedisKey(platform, nodeID))
			}
			_, err = pipe.Exec(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Session) GetPing(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}
