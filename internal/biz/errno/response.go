package errno

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//把错误码转换成返回结构，包括rpc错误码
func ErrorToResponse(ctx context.Context, err error) Response {
	if errCode, ok := err.(*errorCode); ok {
		return errCode.ToResponseWithContext(ctx)
	}
	e := errors.FromError(err)
	rsp := Response{
		Msg: e.Reason,
	}
	code, ok := reasonCodeMap[e.Reason]
	if ok {
		rsp.Code = int(code)
		return rsp
	}
	return ErrUnknown.ToResponseWithContext(ctx)
}
