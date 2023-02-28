package server

import (
	"context"
	"hlayout/internal/server/wire"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/robfig/cron"
)

type CronServer struct {
	cron *cron.Cron
	log  *log.Helper
}

func NewCronServer(services ...wire.CronService) (transport.Server, error) {
	cron := cron.New()
	server := &CronServer{
		cron: cron,
		log:  log.NewHelper(log.GetLogger()),
	}
	for _, svc := range services {
		svc.InitService()
	}
	for _, svc := range services {
		svc.Cron(cron)
	}
	return server, nil
}

func (s *CronServer) Start(ctx context.Context) error {
	s.cron.Start()
	s.log.Infof("cron server start")
	return nil
}

func (s *CronServer) Stop(ctx context.Context) error {
	s.cron.Stop()
	s.log.Infof("cron server stop")
	return nil
}
