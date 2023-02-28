package user

const (
	cstEventTypeLogin = 1 //登录事件
)

type LoginEvent struct {
	UserID     uint64
	ServerID   uint64
	AppID      string
	GameID     uint64
	UserXid    string
	ServerXid  string
	IsNewer    bool
	ChannelExt interface{}
	UnionInt   uint64
	UnionID    string
	Channel    int32
}
