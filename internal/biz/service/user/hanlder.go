package user

import (
	"hlayout/internal/biz/errno"
	"hlayout/internal/conf"
	"hlayout/internal/pkg/sign"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"

	zaplog "gitlab.yeahka.com/gaas/pkg/log"

	jsoniter "github.com/json-iterator/go"

	"hlayout/internal/biz/common"
	aesPkg "hlayout/internal/pkg/aes_pkg"
	"hlayout/internal/pkg/event"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.yeahka.com/gaas/pkg/rpc"
	basepb "gitlab.yeahka.com/gaas/proto/base/v1"
)

func (s *userService) OnUserLogin(ctx *gin.Context) {
	log := zaplog.FromContext(ctx.Request.Context(), "Func", "OnUserLogin")
	//1. 解析参数
	req := userLoginReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("miss param", err.Error())
		ctx.JSON(http.StatusOK, errno.ErrWrongParam.ToResponse())
		return
	}
	env := conf.Env()
	if req.Channel == 0 || (req.GameID == 0 && req.CorpGameID == 0 && req.ServerXID == "") {
		ctx.JSON(http.StatusOK, errno.ErrWrongParam.ToResponse())
		return
	} else if req.Channel == int32(basepb.Channel_Yeahka) && (env == conf.Bootstrap_Online || env == conf.Bootstrap_UAT) {
		ctx.JSON(http.StatusOK, errno.ErrWrongParam.ToResponse())
		return
	}
	//2. 确定服务器信息
	client := basepb.NewBaseClient(rpc.Conn2("base"))
	serverRsp, err := client.GetServerByXid(ctx.Request.Context(), &basepb.GetServerByXidReq{Xid: req.ServerXID, Cgid: req.CorpGameID, Gid: req.GameID})
	if err != nil {
		log.Error("GetServerByXid err:", err.Error())
		ctx.JSON(http.StatusOK, errno.ErrorToResponse(ctx, err))
		return
	}
	brandInfo, err := client.GetBrandInfo(ctx.Request.Context(), &basepb.GetBrandInfoReq{
		BrandId: uint64(serverRsp.BrandId),
	})
	if err != nil {
		log.Error("GetBrandInfo err:", err.Error())
		ctx.JSON(http.StatusOK, errno.ErrorToResponse(ctx, err))
		return
	} else if cgame, ok := brandInfo.Cgames[uint64(serverRsp.Cgid)]; brandInfo.Expired || !ok || !cgame.Opened { //过期品牌走游戏的默认服务器
		serverRsp, err = client.GetServerByXid(ctx.Request.Context(), &basepb.GetServerByXidReq{Gid: serverRsp.Gid})
		if err != nil {
			log.Errorf("GetServerByXid gid:%d default err:%s", serverRsp.Gid, err)
			ctx.JSON(http.StatusOK, errno.ErrorToResponse(ctx, err))
			return
		}
	}
	extend, _ := jsoniter.Marshal(req.Extend.ChannelExt)
	//3. 发起登陆
	loginReq := basepb.LoginReq{
		Channel:    basepb.Channel(req.Channel),
		Sid:        serverRsp.Id,
		Code:       req.Code,
		Ip:         ctx.ClientIP(),
		Extend:     extend,
		DistinctId: req.Extend.DistinctId,
	}
	loginRsp, err := client.Login(ctx.Request.Context(), &loginReq)
	if err != nil {
		log.Errorf("login err:%s", err)
		ctx.JSON(http.StatusOK, errno.ErrorToResponse(ctx, err))
		return
	}
	if loginRsp != nil {
		log.Infof("loginRsp: loginRsp:%#v", *loginRsp)
	}
	//4. 生成token和salt
	token, salt, err := s.genSession(ctx, req.Channel, serverRsp, loginRsp)
	if err != nil {
		log.Errorf("genSession err:%s", err)
		ctx.JSON(http.StatusOK, errno.ErrUnknown.ToResponse())
		return
	}
	//5. 发送登录事件
	err = event.Send(ctx, cstEventTypeLogin, LoginEvent{
		UserID:     loginRsp.Uid,
		ServerID:   serverRsp.Id,
		AppID:      loginRsp.Appid,
		GameID:     uint64(serverRsp.Gid),
		IsNewer:    loginRsp.Newer,
		UserXid:    loginRsp.Uxid,
		ServerXid:  serverRsp.Xid,
		ChannelExt: req.Extend.ChannelExt,
		UnionInt:   loginRsp.UnionInt,
		UnionID:    loginRsp.UnionId,
		Channel:    req.Channel,
	})
	if err != nil {
		log.Errorf("send event err:%s", err)
	}
	//5. 返回对象
	cdata := s.transChannelData(loginRsp.Tag)
	blackUserRsp, err := client.IsBlackUser(ctx.Request.Context(), &basepb.IsBlackUserReq{
		Uid: loginRsp.Uid,
		Sid: serverRsp.Id,
	})
	if err != nil {
		log.Errorf("IsBlackUser uid:%d sid:%d err:%s", loginRsp.Uid, serverRsp.Id, err)
		ctx.JSON(http.StatusOK, errno.ErrorToResponse(ctx, err))
		return
	}
	rsp := userLoginRsp{
		Salt:      salt,
		UserXid:   loginRsp.Uxid,
		UserID:    loginRsp.Uid,
		Token:     token,
		Nick:      loginRsp.Info.Nick,
		Avatar:    loginRsp.Info.Avatar,
		Sex:       int8(loginRsp.Info.Sex),
		ServerXid: serverRsp.Xid,
		GameId:    serverRsp.Gid,
		Newer:     loginRsp.Newer,
		UserIdStr: strconv.FormatUint(loginRsp.Uid, 10),
		CData:     cdata,
		Tags: UserTag{
			Blocking: blackUserRsp.Block,
		},
		Register: loginRsp.Register,
		UnionID:  loginRsp.UnionId,
	}
	ctx.JSON(http.StatusOK, errno.OK.ToDataWithContext(ctx, rsp))
}

func (s *userService) transChannelData(tagJson []byte) map[string]interface{} {
	var tag map[string]interface{}
	if len(tagJson) > 0 {
		if err := util.JSON.Unmarshal(tagJson, &tag); err != nil {
			zaplog.Errorf("unmarshal tag:%s err:%s", string(tagJson), err)
		}
	}
	return tag
}

func (s *userService) genSession(ctx *gin.Context, channel int32, server *basepb.Server, login *basepb.LoginRsp) (string, string, error) {
	token := uuid.New().String()
	saltData := sign.Md5(token + util.SplitUserAgent(ctx.GetHeader("User-Agent")))
	saltKey, err := aesPkg.NewAes(
		aesPkg.WithData(saltData),
		aesPkg.WithIv(login.Uxid[2:18]),
		aesPkg.WithGameKey(login.Gid),
		aesPkg.WithKeyAfterGameKey(login.Appid)).AesEncrypt()
	if err != nil {
		return "", "", err
	}
	version := ctx.GetHeader("X-Gaas-Version")
	session := &common.Session{
		UserID:           login.Uid,
		UserXID:          login.Uxid,
		Salt:             saltData,
		ThirdPartyUserId: login.OpenId, //第三方的id
		Version:          version,
		RegisteTime:      login.RegisteTime,
		UnionInt:         login.UnionInt,
		UnionID:          login.UnionId,
		Channel:          channel,
		Server: common.Server{
			ServerID:   server.Id,
			ServerXID:  server.Xid,
			GameID:     server.Gid,
			CorpGameID: server.Cgid,
			CorpID:     server.CorpId,
			BrandID:    server.BrandId,
			AppID:      login.Appid,
		},
	}
	to := conf.UserServiceCfg().LoginTimeout.AsDuration() + time.Duration(rand.Intn(100))*time.Second
	if err := s.repo.SaveSession(ctx, session, to); err != nil {
		return "", "", err
	} else {
		return token, saltKey, nil
	}
}
