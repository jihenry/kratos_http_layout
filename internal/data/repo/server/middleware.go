package server

import (
	"context"
	"hlayout/internal/biz/common"
	"hlayout/internal/data"
	"hlayout/internal/data/model/cache"
	"hlayout/internal/server/wire"
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"
)

var _ wire.MiddlewareRepo = (*middlewareRepo)(nil)

type middlewareRepo struct {
	data *data.Data
}

func NewMiddlewareRepo(data *data.Data) wire.MiddlewareRepo {
	return &middlewareRepo{data: data}
}

// GetTokenSession implements wire.MiddlewareRepo
func (r *middlewareRepo) GetTokenSession(ctx context.Context, token string) (*common.Session, error) {
	tokenKey := cache.KeyOfLoginSession(token)
	rd, err := r.data.Rdb().Get(ctx, tokenKey).Bytes()
	if err != nil {
		return nil, err
	}
	out := common.Session{}
	if err = util.JSON.Unmarshal(rd, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetUserSalt implements wire.MiddlewareRepo
func (r *middlewareRepo) GetUserSalt(ctx context.Context, userId uint64) (string, error) {
	saltKey := cache.KeyOfSaltKey(userId)
	return r.data.Rdb().Get(ctx, saltKey).Result()
}

// SetSessionExpired implements wire.MiddlewareRepo
func (r *middlewareRepo) SetSessionExpired(ctx context.Context, session *common.Session, to time.Duration) error {
	pipe := r.data.Rdb().Pipeline()
	pipe.Expire(ctx, cache.KeyOfLoginSession(session.Token), to)
	pipe.Expire(ctx, cache.KeyOfSaltKey(session.UserID), to)
	_, err := pipe.Exec(ctx)
	return err
}

// SaveUriIpChain implements wire.MiddlewareRepo
func (r *middlewareRepo) SaveUriIpChain(ctx context.Context, uri string, ipChain string, nonce string, to time.Duration) error {
	cacheKey := uri + "#" + ipChain + "#" + nonce
	if cmd := r.data.Rdb().SetNX(ctx, cacheKey, nonce, to); cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
