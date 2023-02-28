package middleware

import (
	"bytes"
	"hlayout/internal/biz/errno"
	"strconv"
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"

	"github.com/gin-gonic/gin"
	"gitlab.yeahka.com/gaas/pkg/monitor"
)

func Monitor() gin.HandlerFunc {
	return global.monitor()
}

// customResponseWriter 自定义writer
type customResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 写入字节
func (w customResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString 写入字符
func (w customResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (f *middlewareFactory) monitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response errno.Response
		blw := &customResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		start := time.Now()
		c.Next()
		dura := time.Since(start).Milliseconds()
		err := util.JSON.Unmarshal(blw.body.Bytes(), &response)
		if err == nil {
			monitor.Req(c.Request.RequestURI, strconv.FormatInt(int64(response.Code), 10), uint64(dura), map[string]string{
				"msg": response.Msg,
			})
		}
	}
}
