package data

import (
	"context"
	"hlayout/internal/conf"
	"os"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gitlab.yeahka.com/gaas/pkg/cache"
	"gitlab.yeahka.com/gaas/pkg/db"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var ProviderSet = wire.NewSet(NewData)

type Data struct {
	mdb *gorm.DB
	rdb *redis.Client
}

func NewData(c *conf.Data) (*Data, func(), error) {
	dbLogger := db.NewGormLogger(gormlogger.LogLevel(c.Database.LogLevel))
	d := &Data{}
	if c.Database != nil {
		mdb, err := db.NewMysqlClient(
			db.WithSource(c.Database.Source),
			db.WithMaxConn(int(c.Database.MaxConn)),
			db.WithMaxIdleConn(int(c.Database.MaxIdleConn)),
			db.WithMaxLifeTime(c.Database.MaxLifetime.AsDuration()),
			db.WithLogLevel(gormlogger.LogLevel(c.Database.LogLevel)),
			db.WithLogger(dbLogger),
		)
		if err != nil {
			return nil, nil, err
		}
		d.mdb = mdb
	}
	if c.Redis != nil {
		rdb, err := cache.NewRedisClient(
			cache.WithAddr(c.Redis.Addr),
			cache.WithDb(int(c.Redis.Db)),
			cache.WithPassword(c.Redis.Password),
		)
		if err != nil {
			return nil, nil, err
		}
		d.rdb = rdb
	}
	cleanup := func() {
		log.Info("closing the data resources")
		if d.rdb != nil {
			d.rdb.Close()
		}
	}
	return d, cleanup, nil
}

type contextTxKey struct{}

func (d *Data) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.mdb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.mdb
}

func (d *Data) Rdb() *redis.Client {
	return d.rdb
}

func (d *Data) Mdb() *gorm.DB {
	return d.mdb
}

type dbInstance struct {
	Name   string
	Client *gorm.DB
}

type cacheInstance struct {
	Name   string
	Client *redis.Client
}

func (d *Data) AllDB() []*dbInstance {
	out := make([]*dbInstance, 0, 1)
	out = append(out, &dbInstance{
		Name:   "gaas",
		Client: d.mdb,
	})
	return out
}

func (d *Data) AllCache() []*cacheInstance {
	out := make([]*cacheInstance, 0, 1)
	out = append(out, &cacheInstance{
		Name:   "db" + strconv.Itoa(int(d.rdb.Options().DB)),
		Client: d.rdb,
	})
	return out
}

func (d *Data) GetID() string {
	name, _ := os.Hostname() //TODO: 获取实例的唯一ID
	return name
}
