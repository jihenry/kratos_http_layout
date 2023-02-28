package user

type loginExtend struct {
	DistinctId string      `json:"distinctId,omitempty"` //数数游客id
	ChannelExt interface{} `json:"channelExt,omitempty"` //渠道扩展参数
}

type userLoginReq struct {
	Channel    int32       `json:"channel"`          //渠道id，服务器id可以在多个渠道共用
	CorpGameID int32       `json:"corpGameId"`       //企业游戏id，用于那些需要定制渠道参数的商户（比如商户想用自己的微信小游戏主体发布）
	GameID     int32       `json:"gameId"`           //游戏id，用于获取默认参数，因为serverId可能不传
	ServerXID  string      `json:"serverId"`         //服务器id
	Code       string      `json:"code"`             //渠道token
	Extend     loginExtend `json:"extend,omitempty"` //扩展字段
	Sign       string      `json:"sign,omitempty"`   //签名，H5游戏必须签名
}

type UserTag struct {
	Blocking bool `json:"blocking"` //是否已封禁
}

type userLoginRsp struct {
	Salt      string      `json:"salt"`               //盐
	UserXid   string      `json:"userXid"`            //用户xid
	UserID    uint64      `json:"userId"`             //用户id
	Token     string      `json:"token"`              //gaas的token
	Nick      string      `json:"nick,omitempty"`     //昵称
	Avatar    string      `json:"avatar,omitempty"`   //头像
	Sex       int8        `json:"sex,omitempty"`      //性别
	ServerXid string      `json:"serverId,omitempty"` //服务器xid
	GameId    int32       `json:"gameId"`             //游戏id
	Newer     bool        `json:"newer"`              //是否为此服务器下的新手
	UserIdStr string      `json:"userIdStr"`          //用户字符串
	CData     interface{} `json:"cdata,omitempty"`    //渠道数据
	Tags      UserTag     `json:"tags"`               //通用标签数据
	Register  bool        `json:"register"`           //是否刚注册
	UnionID   string      `json:"unionId,omitempty"`  //渠道unionId
}
