package sign

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func RequestIpChain(c *gin.Context) string {
	var (
		buf     strings.Builder
		headers = []string{"X-Forwarded-For", "Proxy-Client-IP", "WL-Proxy-Client-IP", "X-Real-Ip"}
	)
	for _, v := range headers {
		if buf.Len() > 0 {
			buf.WriteByte('|')
		}
		buf.WriteString(v)
		buf.WriteByte('|')
		buf.WriteString(c.Request.Header.Get(v))
	}
	return buf.String()
}
