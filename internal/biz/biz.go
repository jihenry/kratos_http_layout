package biz

import (
	"hlayout/internal/biz/service/cron"
	"hlayout/internal/biz/service/kafka"
	"hlayout/internal/biz/service/user"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	user.NewUserService,
	kafka.NewKafkaCommonConsumer,
	cron.NewCronService,
)
