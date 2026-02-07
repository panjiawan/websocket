package ws

import (
	"context"
	"errors"
	"fmt"
	"github.com/panjiawan/go-lib/pkg/plog"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
	"websocket/conf"
	"websocket/service/dto"
)

type NChanSend struct {
	//wsClient *http.Client
	client *fasthttp.Client
}

func NewNChanSend() *NChanSend {
	return &NChanSend{
		//wsClient: &http.Client{
		//	Transport: &http.Transport{
		//		MaxIdleConns:        200,              //最大空闲连接数
		//		MaxIdleConnsPerHost: 200,              //每个主机的最大空闲连接数
		//		MaxConnsPerHost:     500,              //每个主机的最大连接数
		//		IdleConnTimeout:     30 * time.Second, //空闲连接超时时间
		//	},
		//	Timeout: 5 * time.Second,
		//},
		client: &fasthttp.Client{
			MaxIdleConnDuration: 30 * time.Second, //最大空闲连接时长
			MaxConnsPerHost:     500,              //每个主机的最大连接数
		},
	}
}

const MaxChannel = 250

func (n *NChanSend) SendWebsocket(ctx context.Context, req *dto.WebSocketSendReq) error {
	if len(req.ToUsers) == 0 {
		return nil
	}
	if len(req.ToUsers) <= MaxChannel {
		if err := n.sendMessage(ctx, req.Platform, req.ToUsers, req.Message); err != nil {
			plog.Error("SendWebsocket sendMessage error", zap.Error(err),
				zap.Any("platform", req.Platform), zap.Any("user_ids", req.ToUsers))
			return err
		}
		return nil
	}
	chunks := lo.Chunk(req.ToUsers, MaxChannel)
	var ws sync.WaitGroup
	for _, chunk := range chunks {
		ws.Add(1)
		go func(chunk []string) {
			defer ws.Done()
			if err := n.sendMessage(ctx, req.Platform, chunk, req.Message); err != nil {
				plog.Error("SendWebsocket sendMessage error", zap.Error(err),
					zap.Any("platform", req.Platform), zap.Any("user_ids", chunk))
			}
		}(chunk)
	}
	// 等待所有协程完成
	ws.Wait()
	return nil
}

func (n *NChanSend) sendMessage(ctx context.Context, platform string, userIds []string, message string) error {
	if len(userIds) > MaxChannel {
		return errors.New("max channel count limit 250")
	}

	start := time.Now()
	channels := strings.Join(userIds, ",")
	url := fmt.Sprintf("%s/pub/%s/%s", conf.GetHandle().GetHttpConf().WSHost, platform, channels)
	// 获取请求和响应对象
	httpReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(httpReq)
	httpResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(httpResp)
	// 设置基本请求参数
	httpReq.SetRequestURI(url)
	httpReq.Header.SetMethod(fasthttp.MethodPost)
	//httpReq.Header.SetContentType("application/json")
	httpReq.SetBody([]byte(message))
	httpReq.Header.Set("Content-Type", "application/json")
	// 发送请求
	if err := n.client.DoTimeout(httpReq, httpResp, time.Duration(10)*time.Second); err != nil {
		plog.Error("sendMessage DoTimeout error", zap.Error(err),
			zap.Any("platform", platform), zap.Any("user_ids", channels))
		return err
	}
	defer plog.Debug("sendMessage duration",
		zap.Duration("duration", time.Since(start)),
		zap.String("platform", platform), zap.Int("user_count", len(userIds)))

	// 检查状态码
	statusCode := httpResp.StatusCode()
	if statusCode >= fasthttp.StatusBadRequest {
		plog.Error("sendMessage error", zap.Int("status", statusCode),
			zap.Any("platform", platform), zap.Any("user_ids", channels))
		return errors.New("sendMessage statusCode error")
	}

	return nil
}
