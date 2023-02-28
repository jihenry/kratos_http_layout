package user

import (
	"context"
	"fmt"
	biz_wire "hlayout/internal/biz/wire"
	"hlayout/internal/pkg/event"
	"hlayout/internal/server/wire"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type userService struct {
	log  *log.Helper
	repo biz_wire.UserRepo
}

var _ wire.GinService = (*userService)(nil)

func NewUserService(repo biz_wire.UserRepo) wire.GinService {
	return &userService{log: log.NewHelper(log.GetLogger()), repo: repo}
}

// InitService implements wire.userService
func (s *userService) InitService() error {
	_ = event.Register(event.EventTypeLogin, s.onLoginEvent, &LoginEvent{})
	return nil
}

// UnInitService implements wire.GinService
func (*userService) UnInitService() error {
	return nil
}

func (s *userService) Router(root *gin.RouterGroup) error {
	grg := root.Group("/gaas")
	grg.POST("/user/login", s.OnUserLogin) //用户登录
	return nil
}

func (s *userService) onLoginEvent(ctx context.Context, data interface{}) error {
	eventData, ok := data.(*LoginEvent)
	if !ok {
		err := fmt.Errorf("guest:%v *loginEvent false", data)
		s.log.Error(err)
		return err
	}
	s.log.Infof("onLoginEvent eventData:%v", eventData)
	return nil
}
