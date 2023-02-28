package middleware

import (
	"bytes"
	"hlayout/internal/biz/common"
	"hlayout/internal/biz/errno"
	"hlayout/internal/pkg/sign"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"

	"gitlab.yeahka.com/gaas/pkg/log"

	"github.com/gin-gonic/gin"
)

var (
	nonceTimeOutSeconds = 300
)

func Sign() gin.HandlerFunc {
	return global.sign()
}

func (f *middlewareFactory) sign() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := log.FromContext(c, "Func", "gaasSignVerify")
		saltKey := common.Salt(c)
		if saltKey == "" {
			log.Errorf("salt is empty, userID:%d", common.UserID(c))
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrUnauthorized.ToResponse())
			return
		}
		var (
			byteStr []byte
		)
		if c.Request.Method == "GET" {
			queryMap := make(map[string]interface{})
			data := c.Request.URL.RawQuery
			if strings.Contains(data, "&") {
				vv := strings.Split(data, "&")
				for k := range vv {
					ss := strings.Split(vv[k], "=")
					for i := 0; i < len(ss); i += 2 {
						queryMap[ss[i]] = ss[i+1]
					}
				}
			} else {
				ss := strings.Split(data, "=")
				if len(ss) >= 2 {
					for i := 0; i < len(ss); i += 2 {
						queryMap[ss[i]] = ss[i+1]
					}
				}
			}
			if len(queryMap) > 0 {
				byteStr, _ = util.JSON.Marshal(queryMap)
			} else {
				c.Next()
				return
			}
		} else {
			byteStr, _ = c.Copy().GetRawData()
		}
		csign := c.GetHeader("X-Gaas-Sign")
		if csign == "" {
			log.Errorf("X-Gaas-Sign header miss")
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrInvalidSign.ToResponse())
			return
		}
		ssign, timestamp, nonce := sign.CheckSign(byteStr, saltKey, []string{"timestamp", "nonce"}...)
		if len(ssign) == 0 {
			log.Errorf("server sign empty")
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrInvalidSign.ToResponse())
			return
		}
		for k := range csign {
			if csign[k] != ssign[k] {
				log.Errorf("client sign:%s != server sign:%s sKey:%s", csign, ssign, saltKey)
				c.AbortWithStatusJSON(http.StatusOK, errno.ErrInvalidSign.ToResponse())
				return
			}
		}
		if timestamp == "" {
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrReqRepeat.ToResponse())
			return
		} else {
			cTime, _ := strconv.ParseInt(timestamp, 10, 64)
			//毫秒数
			sTime := time.Now().Unix()
			if sTime-cTime >= int64(nonceTimeOutSeconds) {
				log.Errorf("timeout:%v cTime:%d", sTime-cTime >= int64(nonceTimeOutSeconds), cTime)
				c.AbortWithStatusJSON(http.StatusOK, errno.ErrReqRepeat.ToResponse())
				return
			}
		}
		//判断用户的nonce 参数的有效性
		if nonce == "" {
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrWrongParam.ToResponse())
			return
		}
		ipChain := sign.RequestIpChain(c)
		err := f.repo.SaveUriIpChain(c, c.Request.RequestURI, sign.RequestIpChain(c), nonce, time.Duration(nonceTimeOutSeconds)*time.Second)
		if err != nil {
			log.Errorf("SaveUriIpChain RequestURI%v ipChain:%s nonce:%s err:%s", c.Request.RequestURI, ipChain, nonce, err)
			c.AbortWithStatusJSON(http.StatusOK, errno.ErrReqRepeat.ToResponse())
			return
		}
		if c.Request.Method == "POST" {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(byteStr))
		}
		c.Next()
	}
}
