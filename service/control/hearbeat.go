package control

import (
	"context"
	"github.com/panjiawan/go-lib/pkg/plog"
	"go.uber.org/zap"
	"time"
)

//心跳
func (s *Service) heartbeat() {
	ctx := context.TODO()
	heartbeatTicker := time.NewTicker(time.Minute * 1)
	defer heartbeatTicker.Stop()

	cleanupTicker := time.NewTicker(time.Minute * 5)
	defer cleanupTicker.Stop()

	for {
		select {
		case <-s.exitChn:
			return
		case <-cleanupTicker.C:
			err := s.session.CleanOfflineNodes(ctx)
			if err != nil {
				plog.Error("heartbeat CleanOfflineNodes error", zap.Error(err))
			}
			plog.Debug("heartbeat CleanOfflineNodes Success..")
		case <-heartbeatTicker.C:
			err := s.session.NodeHeartbeat(ctx)
			if err != nil {
				plog.Error("heartbeat NodeHeartbeat error", zap.Error(err))
			}
			plog.Debug("NodeHeartbeat Success..")
		}
	}
}
