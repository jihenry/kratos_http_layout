package server

import (
	"context"
	"errors"
	"fmt"
	"hlayout/internal/conf"
	"hlayout/internal/server/wire"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	monitor "gitlab.yeahka.com/gaas/pkg/monitor"
	"gitlab.yeahka.com/gaas/pkg/util"
)

// monitorServerImpl 监控服务器
type monitorServerImpl struct {
	monitorTime    int
	logfilePath    string
	logger         *log.Helper
	repo           wire.MonitorRepo
	stopCanlelFunc func()
}

// NewMonitorServer 声明一个监控服务器
func NewMonitorServer(repo wire.MonitorRepo) (transport.Server, error) {
	monitorCfg := conf.MonitorCfg()
	if monitorCfg.LogfileName == "" {
		return nil, fmt.Errorf("cfg:%+v is invalid", monitorCfg)
	}
	monitorServer := &monitorServerImpl{
		monitorTime: int(monitorCfg.MonitorTime),
		logfilePath: monitorCfg.LogfileName,
		repo:        repo,
		logger:      log.NewHelper(log.GetLogger()),
	}
	if monitorServer.monitorTime == 0 {
		monitorServer.monitorTime = 10
	}
	return monitorServer, nil
}

// Start 开始监控
func (s *monitorServerImpl) Start(ctx context.Context) error {
	ip, err := util.ExternalIP()
	if err != nil {
		return err
	}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("failed to read build info")
	}
	monitor.Start(ip.String(), bi.Main.Path, s.logfilePath)
	stopCtx, cancel := context.WithCancel(context.Background())
	s.doTickerMonitor(stopCtx)
	s.stopCanlelFunc = cancel
	return nil
}

// Stop implements transport.Server
func (s *monitorServerImpl) Stop(context.Context) error {
	if s.stopCanlelFunc != nil {
		s.stopCanlelFunc()
	}
	monitor.Stop()
	return nil
}

func (s *monitorServerImpl) doTickerMonitor(stopCtx context.Context) {
	t := time.NewTicker(time.Duration(s.monitorTime) * time.Second)
LOOP:
	for {
		select {
		case <-stopCtx.Done():
			break LOOP
		case <-t.C:
			dbMonitorList, err := s.repo.GetDBMonitorData() //db
			if err != nil {
				s.logger.Errorf("GetDBMonitorData err:%s", err)
			}
			for _, v := range dbMonitorList {
				monitor.Src(v.Src, v.Identity, v.Used, v.MaxSrc)
			}
			cacheMonitorList, err := s.repo.GetCacheMonitorData() //缓存
			if err != nil {
				s.logger.Errorf("GetCacheMonitorData err:%s", err)
			}
			for _, v := range cacheMonitorList {
				monitor.Src(v.Src, v.Identity, v.Used, v.MaxSrc)
			}
			threadCount := pprof.Lookup("threadcreate").Count() // 创建的线程数
			gNum := runtime.NumGoroutine()                      // Goroutine数
			monitor.Src("thread_count", "thread", uint64(threadCount), 10000)
			monitor.Src("num_goroutine", "goroutine", uint64(gNum), 10000)
			percent, err := cpu.Percent(time.Second, false) // cpu利用率
			if err != nil {
				continue
			}
			monitor.Src("cpu_percent", "cpu", uint64(percent[0]*float64(100)), 10000)
			memInfo, err := mem.VirtualMemory() // 内存利用率
			if err != nil {
				continue
			}
			monitor.Src("memory_percent", "memory", uint64(memInfo.UsedPercent*float64(100)), 10000)
		}
	}
}
