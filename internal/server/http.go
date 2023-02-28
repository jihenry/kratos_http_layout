package server

import (
	"context"
	"hlayout/internal/conf"
	"hlayout/internal/server/middleware"
	"hlayout/internal/server/wire"

	shttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var _ transport.Server = (*httpServerImpl)(nil)

type httpServerImpl struct {
	server   *http.Server
	services []wire.GinService
	pprof    *shttp.Server
}

func NewHTTPServer(c *conf.Server, middlewareRepo wire.MiddlewareRepo, services ...wire.GinService) transport.Server {
	var opts = []http.ServerOption{}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	for _, svc := range services {
		svc.InitService()
	}
	middleware.InitMiddlewareFactory(middlewareRepo)
	engine := gin.Default()
	engine.Use(
		middleware.Recovery(),
		middleware.Cors(),
		middleware.Request(),
		middleware.Logger(),
		otelgin.Middleware(conf.TraceCfg().ServerName),
		middleware.Monitor(),
	)
	root := engine.Group("/")
	initRoot(root)
	for _, svc := range services {
		svc.Router(root)
	}
	srv := http.NewServer(opts...)
	srv.HandlePrefix("/", engine)
	impl := &httpServerImpl{
		server:   srv,
		services: services,
	}
	if c.Pprof != nil {
		impl.pprof = &shttp.Server{
			Addr:    c.Pprof.Addr,
			Handler: pprof.NewHandler(),
		}
	}
	return impl
}

func (s *httpServerImpl) Start(ctx context.Context) error {
	if s.pprof != nil {
		go s.pprof.ListenAndServe()
	}
	return s.server.Start(ctx)
}

func (s *httpServerImpl) Stop(ctx context.Context) error {
	if s.pprof != nil {
		s.pprof.Shutdown(ctx)
	}
	if err := s.server.Stop(ctx); err != nil {
		return err
	}
	for _, svc := range s.services {
		svc.UnInitService()
	}
	return nil
}

func initRoot(root *gin.RouterGroup) {
	root.GET("/", func(c *gin.Context) {
		c.String(shttp.StatusOK, "ok!")
	})
	root.POST("/ping", func(c *gin.Context) {
		c.String(shttp.StatusOK, "pong")
	})
}
