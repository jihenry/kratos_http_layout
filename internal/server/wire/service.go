package wire

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

type CronService interface {
	InitService() error
	Cron(root *cron.Cron) error
}

type GinService interface {
	InitService() error
	Router(root *gin.RouterGroup) error
	UnInitService() error //服务结束，用于清理资源
}

type KafkaConsumerService interface {
	OnKafkaMessage(ctx context.Context, pm *sarama.ConsumerMessage) error
}
