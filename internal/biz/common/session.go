package common

import (
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Session struct {
	Server
	UserID           uint64 `json:"uid"`         //用户id
	UserXID          string `json:"uxid"`        //用户xid
	Salt             string `json:"salt"`        //当次登陆生成的盐
	ThirdPartyUserId string `json:"tuid"`        //第三方的用户id
	Version          string `json:"version"`     //会话版本，更大范围的版本号，上下文的是某个请求的版本号
	RegisteTime      int64  `json:"registeTime"` //注册时间
	UnionInt         uint64 `json:"unionInt"`    //gaas端的整型unionid
	Channel          int32  `json:"channel"`     //渠道
	UnionID          string `json:"unionId"`     //渠道的unionId
	CreateTime       int64  `json:"createTime"`  //服务器创建时间
	Level            int64  `json:"level"`       //用户等级
	NickName         string `json:"nickName"`    //用户昵称
	UnionServerId    uint64 `json:"-"`
	Token            string `json:"-"`
}

type Server struct {
	GameID     int32  `json:"gid"`     //游戏id
	CorpGameID int32  `json:"cgid"`    //企业游戏id
	ServerID   uint64 `json:"sid"`     //服务器id
	ServerXID  string `json:"sxid"`    //服务器xid
	CorpID     int32  `json:"corpId"`  //企业id
	BrandID    int32  `json:"brandId"` //品牌id
	AppID      string `json:"appId"`   //渠道id
}

const (
	SessionKey = "session"
	VersionKey = "version"
)

func ServerID(ctx context.Context) uint64 {
	session := session(ctx)
	if session != nil {
		return session.ServerID
	}
	return 0
}

func ServerXID(ctx context.Context) string {
	session := session(ctx)
	if session != nil {
		return session.ServerXID
	}
	return ""
}

func Channel(ctx context.Context) int32 {
	session := session(ctx)
	if session != nil {
		return session.Channel
	}
	return 0
}

func UnionInt(ctx context.Context) uint64 {
	session := session(ctx)
	if session != nil {
		return session.UnionInt
	}
	return 0
}

func UnionId(ctx context.Context) string {
	session := session(ctx)
	if session != nil {
		return session.UnionID
	}
	return ""
}

func RegisteTime(ctx context.Context) int64 {
	session := session(ctx)
	if session != nil {
		return session.RegisteTime
	}
	return 0
}

func CreateTime(ctx context.Context) int64 {
	session := session(ctx)
	if session != nil {
		return session.CreateTime
	}
	return 0
}

func Version(ctx context.Context) (bigVersion int, iterVersion int, bugVersion int) {
	version := ctxVersion(ctx)
	if version == "" {
		session := session(ctx)
		if session != nil {
			version = session.Version
		}
	}
	fs := strings.Split(version, ".")
	switch len(fs) {
	case 0:
		return
	case 1:
		bigVersion, _ = strconv.Atoi(fs[0])
	case 2:
		bigVersion, _ = strconv.Atoi(fs[0])
		iterVersion, _ = strconv.Atoi(fs[1])
	default:
		bigVersion, _ = strconv.Atoi(fs[0])
		iterVersion, _ = strconv.Atoi(fs[1])
		bugVersion, _ = strconv.Atoi(fs[2])
	}
	return
}

func ctxVersion(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		v, ok := ginCtx.Get(VersionKey)
		if !ok {
			return ""
		}
		version, _ := v.(string)
		return version
	}
	return getVersion(ctx)
}

func GameID(ctx context.Context) int32 {
	session := session(ctx)
	if session != nil {
		return session.GameID
	}
	return 0
}

func Salt(ctx context.Context) string {
	session := session(ctx)
	if session != nil {
		return session.Salt
	}
	return ""
}

func CorpGameID(ctx context.Context) int32 {
	session := session(ctx)
	if session != nil {
		return session.CorpGameID
	}
	return 0
}

func AllServerInfo(ctx context.Context) Server {
	session := session(ctx)
	if session != nil {
		return session.Server
	}
	return Server{}
}

func UserID(ctx context.Context) uint64 {
	session := session(ctx)
	if session != nil {
		return session.UserID
	}
	return 0
}

func UserXID(ctx context.Context) string {
	session := session(ctx)
	if session != nil {
		return session.UserXID
	}
	return ""
}

func AppID(ctx context.Context) string {
	session := session(ctx)
	if session != nil {
		return session.AppID
	}
	return ""
}

func CorpID(ctx context.Context) int32 {
	session := session(ctx)
	if session != nil {
		return session.CorpID
	}
	return 0
}

func BrandID(ctx context.Context) int32 {
	session := session(ctx)
	if session != nil {
		return session.BrandID
	}
	return 0
}

// ThirdPartyUserId 获取第三方用户Id
//func ThirdPartyUserId(ctx context.Context) string {
//	session := session(ctx)
//	if session != nil {
//		return session.ThirdPartyUserId
//	}
//	return ""
//}

func session(ctx context.Context) *Session {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		v, ok := ginCtx.Get(SessionKey)
		if !ok {
			return nil
		}
		session, _ := v.(*Session)
		return session
	}
	return getSession(ctx)
}

type requestSession struct {
}
type requestVersion struct {
}

func SetSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, requestSession{}, session)
}

func getSession(ctx context.Context) *Session {
	if val, ok := ctx.Value(requestSession{}).(*Session); ok {
		return val
	}
	return nil
}

func getVersion(ctx context.Context) string {
	if val, ok := ctx.Value(requestVersion{}).(string); ok {
		return val
	}
	return ""
}

func SetVersion(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, requestVersion{}, version)
}

// UserLevel 用户等级
func UserLevel(ctx context.Context) int64 {
	if ss := session(ctx); ss != nil {
		return ss.Level
	}
	return 0
}

func UserNick(ctx context.Context) string {
	if ss := session(ctx); ss != nil {
		return ss.NickName
	}
	return ""
}

func UnionServerId(ctx context.Context) uint64 {
	if ss := session(ctx); ss != nil {
		return ss.UnionServerId
	}
	return uint64(0)
}
