package event

import (
	"context"
	"fmt"
	"hlayout/internal/conf"
	"reflect"
	"strconv"
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/Shopify/sarama"
	"gitlab.yeahka.com/gaas/pkg/mq/kafka"
)

type EventType int32

type EventData struct {
	Type  EventType //事件类型
	ID    int32     //事件id
	Time  int64     //任务发送时间
	Data  []byte    //携带数据
	Topic string    `json:"-"`
}

var pkEventConfig = map[EventType]interface{}{}

var eventFuncData = map[EventType]eventFunc{}

type eventFunc func(ctx context.Context, data interface{}) error

func Register(event EventType, cb eventFunc, dataType interface{}) error {
	f, ok := eventFuncData[event]
	if ok {
		return fmt.Errorf("the func exist:%#v", f)
	}
	eventFuncData[event] = cb
	pkEventConfig[event] = dataType
	return nil
}

type SendOption func(*EventData)

func WithTopic(topic string) SendOption {
	return func(te *EventData) {
		te.Topic = topic
	}
}

func Send(ctx context.Context, etype EventType, data interface{}, opts ...SendOption) error {
	kafkaProducerCfg := conf.KafkaProducerCfg()
	log.Infof("Send etype:%v data:%#v", etype, data)
	event := EventData{
		Type:  etype,
		Time:  time.Now().Unix(),
		Topic: kafkaProducerCfg.Common.Topic,
	}
	for _, opt := range opts {
		opt(&event)
	}
	if do, err := util.JSON.Marshal(data); err != nil {
		return err
	} else {
		event.Data = do
	}
	mo, err := util.JSON.Marshal(event)
	if err != nil {
		return err
	}
	err = kafka.SendAsyncMsg(ctx, &sarama.ProducerMessage{
		Topic: event.Topic,
		Key:   sarama.ByteEncoder(strconv.Itoa(int(etype))),
		Value: sarama.ByteEncoder(mo),
	})
	if err != nil {
		log.Errorf("Send etype:%d push fail:%s", etype, err)
	}
	return err
}

func Dispatch(ctx context.Context, event EventType, data []byte) error {
	dataTmpl, ok := pkEventConfig[event]
	if !ok {
		return fmt.Errorf("the event:%d is unknown", event)
	}
	f, ok := eventFuncData[event]
	if !ok {
		return fmt.Errorf("event:%d deal func not find", event)
	}
	dataType := reflect.TypeOf(dataTmpl)
	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
	}
	dispathData := reflect.New(dataType)
	err := util.JSON.Unmarshal(data, dispathData.Interface())
	if err != nil {
		return err
	}
	return f(context.Background(), dispathData.Interface())
}
