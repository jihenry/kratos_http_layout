package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"runtime"

	"gitlab.yeahka.com/gaas/pkg/log"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return global.recovery()
}

func (f *middlewareFactory) recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				data, _ := c.GetRawData()
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				log.FromContext(c).Errorf("%v: %+v\n%s\n", err, string(data), buf)
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				c.AbortWithStatus(200)
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": http.StatusBadGateway,
					"msg":  http.StatusText(http.StatusBadGateway),
					"data": nil,
				})
				return
			}
		}()
		c.Next()
	}
}
