package kafka

import (
	"context"
	"hlayout/internal/pkg/event"
	"hlayout/internal/server/wire"
	"runtime/debug"
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"

	"github.com/Shopify/sarama"
	"github.com/go-kratos/kratos/v2/log"
)

var _ wire.KafkaConsumerService = (*commonKafkaService)(nil)

type commonKafkaService struct {
	log *log.Helper
}

func NewKafkaCommonConsumer() wire.KafkaConsumerService {
	return &commonKafkaService{
		log: log.NewHelper(log.GetLogger()),
	}
}

// OnKafkaMessage implements wire.KafkaConsumerService
func (s *commonKafkaService) OnKafkaMessage(ctx context.Context, pm *sarama.ConsumerMessage) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("recover:%s %v", err, string(debug.Stack()))
		}
	}()
	var eventData event.EventData
	err := util.JSON.Unmarshal(pm.Value, &eventData)
	if err != nil {
		return err
	}
	log.Infof("recv eventType:%d id:%d data:%s sendTime:%s", eventData.Type, eventData.ID,
		string(eventData.Data), time.Unix(eventData.Time, 0))
	err = event.Dispatch(ctx, eventData.Type, eventData.Data)
	if err != nil {
		s.log.Errorf("dispatchEvent eventType:%d id:%d presult:%#v err:%s", eventData.Type, eventData.ID, eventData.Data, err)
		return err
	}
	return nil
}
