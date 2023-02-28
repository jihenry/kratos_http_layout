package wire

import (
	"context"
	"hlayout/internal/biz/common"
	"time"
)

type UserRepo interface {
	SaveSession(ctx context.Context, session *common.Session, to time.Duration) error
}
