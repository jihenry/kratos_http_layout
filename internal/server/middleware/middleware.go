package middleware

import (
	"hlayout/internal/server/wire"
)

var global *middlewareFactory

type middlewareFactory struct {
	repo wire.MiddlewareRepo
}

func InitMiddlewareFactory(repo wire.MiddlewareRepo) {
	global = &middlewareFactory{repo: repo}
}
