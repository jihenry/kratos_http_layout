package errno

import (
	"fmt"

	activitypb "gitlab.yeahka.com/gaas/proto/activity/v1"
	basepb "gitlab.yeahka.com/gaas/proto/base/v1"
	couponpb "gitlab.yeahka.com/gaas/proto/coupon/v1"
	developbasepb "gitlab.yeahka.com/gaas/proto/develop/base/v1"
	blendpb "gitlab.yeahka.com/gaas/proto/develop/blend/v1"
	dynamicpb "gitlab.yeahka.com/gaas/proto/develop/dynamic/v1"
	publicpb "gitlab.yeahka.com/gaas/proto/develop/public/v1"
	errorpb "gitlab.yeahka.com/gaas/proto/error/v1"
	gamepb "gitlab.yeahka.com/gaas/proto/game/v1"
	itempb "gitlab.yeahka.com/gaas/proto/item/v1"
	nprizepb "gitlab.yeahka.com/gaas/proto/nprize/v1"
	prizepb "gitlab.yeahka.com/gaas/proto/prize/v1"
	rankpb "gitlab.yeahka.com/gaas/proto/rank/v1"
	tactpb "gitlab.yeahka.com/gaas/proto/tact/v1"
	taskpb "gitlab.yeahka.com/gaas/proto/task/v1"
)

var reasonCodeMap = map[string]int32{}
var codeReasonMap = map[int32]string{}

func init() {
	buildPBErrNoMap(activitypb.ErrorReason_value)
	buildPBErrNoMap(couponpb.ErrorReason_value)
	buildPBErrNoMap(prizepb.ErrorReason_value)
	buildPBErrNoMap(rankpb.ErrorReason_value)
	buildPBErrNoMap(tactpb.ErrorReason_value)
	buildPBErrNoMap(nprizepb.ErrorReason_value)
	buildPBErrNoMap(itempb.ErrorReason_value)
	buildPBErrNoMap(gamepb.ErrorReason_value)
	buildPBErrNoMap(taskpb.ErrorReason_value)
	buildPBErrNoMap(basepb.ErrorReason_value)
	buildPBErrNoMap(blendpb.ErrorReason_value)
	buildPBErrNoMap(dynamicpb.ErrorReason_value)
	buildPBErrNoMap(publicpb.ErrorReason_value)
	buildPBErrNoMap(developbasepb.ErrorReason_value)
	buildPBErrNoMap(errorpb.ErrorReason_value)
}

func buildPBErrNoMap(pbReasonValue map[string]int32) {
	for reason, code := range pbReasonValue {
		if code <= 0 {
			continue
		}
		if existReason, ok := codeReasonMap[code]; ok {
			panic(fmt.Sprintf("reason:%s with %s code:%d is repeatd", reason, existReason, code))
		}
		reasonCodeMap[reason] = code
		codeReasonMap[code] = reason
	}
}
