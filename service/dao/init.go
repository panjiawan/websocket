package dao

import (
	"fmt"
	"github.com/panjiawan/go-lib/pkg/plog"
	"github.com/panjiawan/go-lib/pkg/predis"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
	"websocket/conf"
)

var (
	RedisKey = "default"
)

func Run() {
	InitRedis()
}

var (
	redisHandles   map[string]*predis.Service
	redisKeyPrefix map[string]string
)

func InitRedis() {
	redisConf := conf.GetHandle().GetRedisConf()

	redisHandles = make(map[string]*predis.Service)
	redisKeyPrefix = make(map[string]string)
	if len(redisConf.Hosts) == 1 {
		RedisKey = redisConf.Hosts[0].Name
	}

	for _, cfg := range redisConf.Hosts {
		timeout := time.Duration(cfg.Timeout) * time.Second
		redisHandles[cfg.Name] = predis.New(
			predis.WithConnection(cfg.Host, cfg.Port),
			predis.WithAuth(cfg.Auth),
			predis.WithDB(cfg.DB),
			predis.WithLimit(cfg.MinIdle, cfg.MaxIdle),
			predis.WithReadTimeout(timeout),
			predis.WithWriteTimeout(timeout),
		)

		redisKeyPrefix[cfg.Name] = cfg.Prefix
	}

	for k, f := range redisHandles {
		if err := f.Run(); err != nil {
			plog.Error("redis start error", zap.String("key", k), zap.Error(err))
			panic(err)
		} else {
			plog.Info("redis started", zap.String("key", k))
		}
	}
}

func Redis(key ...string) *redis.Client {
	if len(key) == 1 {
		return redisHandles[key[0]].GetConn()
	}
	return redisHandles[RedisKey].GetConn()
}

func FormatRedisKey(key string) string {
	return fmt.Sprintf("%s:%s", redisKeyPrefix[RedisKey], key)
}
