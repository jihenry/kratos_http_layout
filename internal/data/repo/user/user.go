package user

import (
	"context"
	"hlayout/internal/biz/common"
	"hlayout/internal/biz/wire"
	"hlayout/internal/data"
	"hlayout/internal/data/model/cache"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var _ wire.UserRepo = (*userRepoImpl)(nil)

type userRepoImpl struct {
	log  *log.Helper
	data *data.Data
}

func NewUserRepo(data *data.Data) wire.UserRepo {
	return &userRepoImpl{data: data, log: log.NewHelper(log.GetLogger())}
}

// Login implements wire.UserRepo
func (r *userRepoImpl) SaveSession(ctx context.Context, session *common.Session, to time.Duration) error {
	pipe := r.data.Rdb().Pipeline()
	pipe.Expire(ctx, cache.KeyOfLoginSession(session.Token), to)
	pipe.Expire(ctx, cache.KeyOfSaltKey(session.UserID), to)
	_, err := pipe.Exec(ctx)
	return err
}
