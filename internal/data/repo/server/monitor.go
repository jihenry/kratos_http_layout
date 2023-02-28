package server

import (
	"hlayout/internal/data"
	"hlayout/internal/server/wire"
)

var _ wire.MonitorRepo = (*monitorRepoImpl)(nil)

type monitorRepoImpl struct {
	data *data.Data
}

func NewMonitorRepo(data *data.Data) wire.MonitorRepo {
	return &monitorRepoImpl{data: data}
}

// GetCacheMonitorData implements wire.MonitorRepo
func (r *monitorRepoImpl) GetCacheMonitorData() ([]wire.MonitorSrc, error) {
	all := r.data.AllCache()
	out := make([]wire.MonitorSrc, 0, len(all))
	for _, v := range all {
		if v == nil {
			continue
		}
		stats := v.Client.PoolStats()
		out = append(out, wire.MonitorSrc{
			Src:      "redis_pool_usage",
			Identity: v.Name,
			Used:     uint64(stats.TotalConns - stats.IdleConns),
			MaxSrc:   uint64(stats.TotalConns),
		})
	}
	return out, nil
}

// GetDBMonitorData implements wire.MonitorRepo
func (r *monitorRepoImpl) GetDBMonitorData() ([]wire.MonitorSrc, error) {
	all := r.data.AllDB()
	out := make([]wire.MonitorSrc, 0, len(all))
	for _, v := range all {
		if v == nil {
			continue
		}
		sqlDb, err := v.Client.DB()
		if err != nil {
			continue
		}
		out = append(out, wire.MonitorSrc{
			Src:      "mysql_pool_usage",
			Identity: v.Name,
			Used:     uint64(sqlDb.Stats().InUse),
			MaxSrc:   uint64(sqlDb.Stats().OpenConnections),
		})
	}
	return out, nil
}
