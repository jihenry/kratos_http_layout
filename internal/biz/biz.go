package biz

import (
	"hlayout/internal/biz/service/activity"
	"hlayout/internal/biz/service/base"
	"hlayout/internal/biz/service/game"
	"hlayout/internal/biz/service/game/crush"
	"hlayout/internal/biz/service/game/dfgirl"
	"hlayout/internal/biz/service/game/pingtu"
	"hlayout/internal/biz/service/game/pocket"
	"hlayout/internal/biz/service/game/watermelon"
	"hlayout/internal/biz/service/group"
	"hlayout/internal/biz/service/plugin"
	"hlayout/internal/biz/service/prize"
	"hlayout/internal/biz/service/report"
	"hlayout/internal/biz/service/share"
	"hlayout/internal/biz/service/task"
	"hlayout/internal/biz/service/user"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	base.NewBaseService,
	activity.NewActivityService,
	watermelon.NewWatermelonService,
	crush.NewCrushService,
	dfgirl.NewdfgirlService,
	share.NewShareService,
	prize.NewPrizeService,
	game.NewGameService,
	pocket.NewPocketService,
	plugin.NewTreeService,
	group.NewGroupService,
	task.NewTaskService,
	report.NewReportService,
	plugin.NewPluginService,
	user.NewUserService,
	pingtu.NewPingtuService,
)
