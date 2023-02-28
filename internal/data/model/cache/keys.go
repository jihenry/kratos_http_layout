package cache

import (
	"strconv"
)

func KeyOfLoginSession(token string) string {
	return "gaas:token:" + token
}

// KeyOfSaltKey  用户的登录生成的签名盐
func KeyOfSaltKey(userID uint64) string {
	return "gaas:salt:" + strconv.FormatUint(userID, 10)
}

func KeyOfEventKey(etype int32) string {
	return "gaas:event:" + strconv.FormatInt(int64(etype), 10)
}
