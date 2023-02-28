package errno

import (
	"context"
	"fmt"
	"hlayout/internal/biz/common"
)

type errorCode struct {
	Code  int
	Msg   string
	Desc  string
	NCode int //新错误码
}

//new 错误码构建函数
func new(msg string, code int, desc string) *errorCode {
	return &errorCode{Msg: msg, Code: code, Desc: desc}
}

//newOldError 给旧错误码迁移用的初始化函数
func newOldError(msg string, code int, ncode int, desc string) *errorCode {
	return &errorCode{Msg: msg, Code: code, Desc: desc, NCode: ncode}
}

func (r *errorCode) Error() string { return fmt.Sprintf("msg:%s desc:%s", r.Msg, r.Desc) }

//WithMsg 定制msg
func (r *errorCode) WithMsg(msg string) *errorCode {
	ec := *r
	ec.Msg = msg
	return &ec
}

//ToResponse 把错误码转换成返回结构
func (r *errorCode) ToResponse() Response {
	return Response{
		Code: r.Code,
		Msg:  r.Msg,
	}
}

func (r *errorCode) ToMessage() Response {
	return Response{
		Code: r.Code,
		Msg:  r.Desc,
	}
}

//ToData 把错误码转换成携带数据的返回结构
func (r *errorCode) ToData(data interface{}) Response {
	rsp := Response{
		Code: r.Code,
		Msg:  r.Msg,
	}
	if data != nil {
		rsp.Data = data
	}
	return rsp
}

//ToResponseWithContext 把错误码转换成返回结构，增加上下文逻辑
func (r *errorCode) ToResponseWithContext(ctx context.Context) Response {
	code := r.Code
	if r.NCode > 0 && isUseNewCode(ctx) {
		code = r.NCode
	}
	return Response{
		Code: code,
		Msg:  r.Msg,
	}
}

//ToDataWithContext 把错误码转换成携带数据的返回结构，增加上下文逻辑
func (r *errorCode) ToDataWithContext(ctx context.Context, data interface{}) Response {
	code := r.Code
	if r.NCode > 0 && isUseNewCode(ctx) {
		code = r.NCode
	}
	rsp := Response{
		Code: code,
		Msg:  r.Msg,
	}
	if data != nil {
		rsp.Data = data
	}
	return rsp
}

func isUseNewCode(ctx context.Context) bool {
	bigVersion, iterVersion, _ := common.Version(ctx)
	gameId := common.GameID(ctx)
	switch {
	case common.IsBoardGame(uint64(gameId)): //神奇板子>=1.8.0
		if bigVersion > 1 {
			return true
		} else if bigVersion == 1 {
			return iterVersion >= 8
		} else {
			return false
		}
	default:
		return false
	}
}
