package main

import (
	"hlayout/internal/biz/service/cron"
	kafkaConsumer "hlayout/internal/biz/service/kafka"
	"hlayout/internal/biz/service/user"
	"hlayout/internal/conf"
	"hlayout/internal/data"
	rserver "hlayout/internal/data/repo/server"
	ruser "hlayout/internal/data/repo/user"
	"hlayout/internal/server"
	"hlayout/internal/server/wire"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
)

//TODO: 目前手动写的，使用wire自动生成
func wireApp(confData *conf.Data, confServer *conf.Server, registry *nacos.Registry) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData)
	if err != nil {
		return nil, nil, err
	}
	middlewareRepo := rserver.NewMiddlewareRepo(dataData)
	userRepo := ruser.NewUserRepo(dataData)
	monitorRepo := rserver.NewMonitorRepo(dataData)
	var ginServices []wire.GinService
	//插件服务
	ginServices = append(ginServices, user.NewUserService(userRepo))
	//定时服务
	cronService := cron.NewCronService()
	//kafka服务
	commonKafkaConsumer := kafkaConsumer.NewKafkaCommonConsumer()
	httpServer := server.NewHTTPServer(confServer, middlewareRepo, ginServices...)
	cronServer, err := server.NewCronServer(cronService)
	if err != nil {
		return nil, nil, err
	}
	monitorServer, err := server.NewMonitorServer(monitorRepo)
	if err != nil {
		return nil, nil, err
	}
	kafkaConsumerServer, err := server.NewKafkaConsumeServer(commonKafkaConsumer)
	if err != nil {
		return nil, nil, err
	}
	app := newApp(registry, httpServer, cronServer, monitorServer, kafkaConsumerServer)
	return app, func() {
		cleanup()
	}, nil
}
