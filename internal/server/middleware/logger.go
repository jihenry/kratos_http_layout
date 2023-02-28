package middleware

import (
	"bytes"
	"hlayout/internal/biz/common"
	"hlayout/internal/biz/errno"
	"hlayout/internal/conf"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"gitlab.yeahka.com/gaas/pkg/log"
	"gitlab.yeahka.com/gaas/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

var (
	reqExcludeHeaderDump = map[string]bool{
		"Host":              true,
		"Transfer-Encoding": false,
		"Trailer":           false,
		"Accept":            false,
		"Accept-Encoding":   false,
		"Connection":        false,
		"Cache-Control":     true,
		"Accept-Language":   true,
		"Origin":            true,
		"Sec-Fetch-Site":    false,
	}
)

type loggerParam struct {
	Start      time.Time
	ClientIP   string
	Req        *dumpRequestDetail
	RespStatus int
	Latency    time.Duration
	RespBody   string
	WxUserId   uint64
	ServerId   uint64
}

func Logger() gin.HandlerFunc {
	return global.logger()
}

func (f *middlewareFactory) logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		dumpReq, err := dumpRequest(c.Request, reqExcludeHeaderDump)
		if err != nil {
			log.Errorf("mw.OperationLog dump req failed. err=%v", err)
		}
		blw := &customResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		// Stop timer
		latency := time.Since(start)
		if latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			latency = latency - latency%time.Second
		}
		wxUserIntId := common.UserID(c)
		serverIntId := common.ServerID(c)
		lp := loggerParam{
			Start:      start,
			ClientIP:   c.ClientIP(),
			RespStatus: c.Writer.Status(),
			Latency:    latency,
			Req:        dumpReq,
			WxUserId:   wxUserIntId,
			ServerId:   serverIntId,
			RespBody:   "",
		}
		if conf.MiddlewareCfg().LogReply {
			lp.RespBody = blw.body.String()
		} else {
			rsp := errno.Response{}
			if err := util.JSON.Unmarshal(blw.body.Bytes(), &rsp); err != nil {
				lp.RespBody = "jsonIterator.Unmarshal err:" + err.Error()
			} else if rsp.Code != 0 {
				lp.RespBody = blw.body.String()
			}
		}
		log.FromContext(c.Request.Context()).Infof("[GIN] ip:%s | method:%s | uri:%s | headers:%s | body:'%s' | latency:%s | status:%d | wxuserid:%d | serverid:%d |trace_id:%v|span_id:%v|resp:%s",
			lp.ClientIP,
			lp.Req.Method,
			lp.Req.URI,
			lp.Req.Headers,
			lp.Req.Body,
			lp.Latency,
			lp.RespStatus,
			lp.WxUserId,
			lp.ServerId,
			tracing.TraceID()(c.Request.Context()),
			tracing.SpanID()(c.Request.Context()),
			lp.RespBody,
		)
	}
}

type dumpRequestDetail struct {
	Method  string      // method
	Path    string      // req path eg: /api/user
	URI     string      // req uri eg: /api/user?id=1
	Body    []byte      // req body
	Headers http.Header // req header
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func dumpRequest(req *http.Request, headDump map[string]bool) (*dumpRequestDetail, error) {
	reqURI := req.RequestURI
	if reqURI == "" {
		reqURI = req.URL.RequestURI()
	}
	detail := dumpRequestDetail{
		Method:  req.Method,
		Path:    req.URL.Path,
		URI:     reqURI,
		Headers: make(http.Header, 0),
	}
	for k, v := range req.Header {
		if h, ok := headDump[k]; !(h && ok) {
			detail.Headers[k] = v
		}
	}
	var err error
	save := req.Body
	defer func() {
		req.Body = save
	}()
	save, req.Body, err = drainBody(req.Body)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	if req.Body != nil {
		var dest io.Writer = &b
		chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
		if chunked {
			dest = httputil.NewChunkedWriter(dest)
		}
		_, err = io.Copy(dest, req.Body)
		if chunked {
			_ = dest.(io.Closer).Close()
		}
		if err != nil {
			return nil, err
		}
	}
	detail.Body = b.Bytes()
	return &detail, err
}
