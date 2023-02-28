package middleware

import (
	"hlayout/internal/biz/common"
	"hlayout/internal/biz/errno"
	"hlayout/internal/conf"
	"math/rand"
	"net/http"
	"time"

	plog "gitlab.yeahka.com/gaas/pkg/log"

	"github.com/gin-gonic/gin"
)

func Token() gin.HandlerFunc {
	return global.token()
}

func (f *middlewareFactory) token() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := plog.FromContext(c, "Middleware", "Token")
		token := c.GetHeader("X-Gaas-Token")
		if token == "" {
			log.Errorf("X-Gaas-Token header miss")
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrUnauthorized.ToResponse())
			return
		}
		session, err := f.repo.GetTokenSession(c, token)
		if err != nil {
			log.Error("GetTokenSession error:%s", err)
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrUnauthorized.ToResponse())
			return
		}
		session.Token = token
		curSalt, err := f.repo.GetUserSalt(c, session.UserID)
		if err != nil {
			log.Error("GetUserSalt error:%s", err)
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrUnauthorized.ToResponse())
			return
		}
		if curSalt != session.Salt {
			log.Errorf("session timeout userID:%d curSalt:%s != session.Salt:%s", session.UserID, curSalt, session.Salt)
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrUnauthorized.ToResponse())
			return
		}
		//TODO: 得设置个最大值
		timeout := conf.UserServiceCfg().LoginTimeout.AsDuration() + time.Duration(rand.Intn(1000)*int(time.Millisecond))
		if err := f.repo.SetSessionExpired(c, session, timeout); err != nil {
			log.Errorf("SetSessionExpired err:%s", err)
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrUnauthorized.ToResponse())
			return
		}
		version := c.GetHeader("X-Gaas-Version")
		c.Set(common.SessionKey, session)
		c.Set(common.VersionKey, version) //每个请求可能用不同的版本，不放在session里面
		//设置用户基础信息到context.Context 上下文中
		requestCtx := common.SetSession(c.Request.Context(), session)
		requestCtx = common.SetVersion(requestCtx, version)
		requestCtx = plog.NewContext(requestCtx, plog.FromLogger(requestCtx,
			`Request-User-Id`, session.UserID,
			`Request-Token`, token,
			`Request-Server-Id`, session.ServerID))
		c.Request = c.Request.WithContext(requestCtx)
		c.Next()
	}
}
