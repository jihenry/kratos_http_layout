package repo

import (
	"hlayout/internal/data/repo/server"
	"hlayout/internal/data/repo/user"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(server.NewMiddlewareRepo, user.NewUserRepo)
