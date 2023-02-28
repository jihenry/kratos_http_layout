package middleware

import (
	"gitlab.yeahka.com/gaas/pkg/log"
	"gitlab.yeahka.com/gaas/pkg/middleware"
	"gitlab.yeahka.com/gaas/pkg/util"

	"github.com/gin-gonic/gin"
)

func Request() gin.HandlerFunc {
	return global.request()
}

func (f *middlewareFactory) request() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := util.New()
		requestDataCtx := middleware.NewRequestDataContext(c.Request.Context(), &middleware.RequestData{RequestId: requestId})
		loggerCtx := log.NewContext(requestDataCtx, log.With(log.GetLogger(), middleware.RequestKeyXRequestID, requestId))
		c.Request = c.Request.WithContext(loggerCtx) //设置到请求体中
		c.Next()
	}
}
