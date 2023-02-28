package cron

import (
	"hlayout/internal/conf"
	"hlayout/internal/server/wire"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron"
)

type cronUsecase struct {
	log *log.Helper
}

var _ wire.CronService = (*cronUsecase)(nil)

func NewCronService() wire.CronService {
	return &cronUsecase{log: log.NewHelper(log.GetLogger())}
}

// InitService implements wire.CronService
func (*cronUsecase) InitService() error {
	return nil
}

// Cron implements server.CronService
func (c *cronUsecase) Cron(root *cron.Cron) error {
	subCron := conf.CronCfg().Sub
	err := root.AddFunc(subCron, func() {
		c.log.Infof("sub cron start")
	})
	if err != nil {
		c.log.Error(err)
		return err
	}

	return nil
}
