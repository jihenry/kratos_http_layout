package wire

import (
	"context"
	"hlayout/internal/biz/common"
	"time"
)

type MiddlewareRepo interface {
	GetTokenSession(ctx context.Context, token string) (*common.Session, error)
	GetUserSalt(ctx context.Context, userId uint64) (string, error)
	SetSessionExpired(ctx context.Context, session *common.Session, to time.Duration) error
	SaveUriIpChain(ctx context.Context, uri, ipChain string, nonce string, to time.Duration) error
}

type MonitorSrc struct {
	Src      string
	Identity string
	Used     uint64
	MaxSrc   uint64
}

type MonitorRepo interface {
	GetDBMonitorData() ([]MonitorSrc, error)
	GetCacheMonitorData() ([]MonitorSrc, error)
}
