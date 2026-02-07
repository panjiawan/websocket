package control

import (
	"context"
	"websocket/service/control/consumer"
	"websocket/service/control/producer"
	"websocket/service/control/session"
	"websocket/service/defines"
)

type Service struct {
	producer producer.IProducer
	consumer consumer.IConsumer
	session  *session.Session
	exitChn  chan struct{}
}

var servicer *Service

func Run() {
	newSession := session.NewSession()

	servicer = &Service{
		producer: producer.NewRedisProducer(),
		consumer: consumer.NewRedisConsumer(newSession),
		session:  newSession,
		exitChn:  make(chan struct{}),
	}

	// 启动一个心跳线程
	go servicer.heartbeat()

	for _, v := range defines.PlatformAll {
		go servicer.consumer.ConsumeMessage(context.TODO(), v)
	}
}

func Stop() {
	servicer.exitChn <- struct{}{}
	servicer.session.CleanSelfSession(context.TODO())
	close(servicer.exitChn)
}
