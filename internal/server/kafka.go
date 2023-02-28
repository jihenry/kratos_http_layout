package server

import (
	"context"
	"hlayout/internal/conf"
	"hlayout/internal/server/wire"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/go-kratos/kratos/v2/transport"
	"gitlab.yeahka.com/gaas/pkg/mq/kafka"
	"gitlab.yeahka.com/gaas/pkg/util"
)

type consumerServiceItem struct {
	consumer kafka.KafkaConsumer
	service  wire.KafkaConsumerService
}

type kafkaConsumerServerImpl struct {
	consumerServices map[string]consumerServiceItem
	stopFunc         context.CancelFunc
	timeout          time.Duration
}

func NewKafkaConsumeServer(commonService wire.KafkaConsumerService) (transport.Server, error) {
	kafkaCfg := conf.KafkaCfg()
	kafkaConsumerCfg := conf.KafkaConsumerCfg()
	common, err := kafka.NewKafkaConsumer(kafkaCfg.Addrs, []string{kafkaConsumerCfg.Common.Topic}, kafkaConsumerCfg.Common.Group,
		kafka.WithConsumeMode(cluster.ConsumerMode(kafkaConsumerCfg.Common.Mode)), kafka.WithConsumeName("common"))
	if err != nil {
		return nil, err
	}
	consumerServices := make(map[string]consumerServiceItem, 1)
	consumerServices["common"] = consumerServiceItem{
		consumer: common,
		service:  commonService,
	}
	impl := &kafkaConsumerServerImpl{consumerServices: consumerServices}
	impl.timeout = kafkaConsumerCfg.Timeout.AsDuration()
	return impl, nil
}

// Start implements transport.Server
func (s *kafkaConsumerServerImpl) Start(ctx context.Context) error {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	for _, consumerService := range s.consumerServices {
		go consumerService.consumer.Receive(cancelCtx, s.onKafkaMessage(cancelCtx, consumerService.service))
	}
	s.stopFunc = cancelFunc
	return nil
}

func (s *kafkaConsumerServerImpl) onKafkaMessage(ctx context.Context, service wire.KafkaConsumerService) func(pm *sarama.ConsumerMessage) error {
	return func(pm *sarama.ConsumerMessage) error {
		defer util.Recover("kafkaConsumer")
		stopCtx, cancel := context.WithTimeout(ctx, s.timeout)
		defer cancel()
		return service.OnKafkaMessage(stopCtx, pm)
	}
}

// Stop implements transport.Server
func (s *kafkaConsumerServerImpl) Stop(ctx context.Context) error {
	s.stopFunc()
	for _, consumerService := range s.consumerServices {
		if err := consumerService.consumer.Stop(); err != nil {
			return err
		}
	}
	return nil
}
