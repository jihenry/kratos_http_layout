package common

import (
	"hlayout/internal/conf"
)

//水果合合合
func GetGameSGHHHID() int32 {
	return 1
}

//消消乐
func GetGameCrushID() int32 {
	switch conf.Env() {
	case conf.Bootstrap_UAT, conf.Bootstrap_Online:
		return 2
	default:
		return 3
	}
}

//悦动女孩
func GetGameDfgirlID() int32 {
	switch conf.Env() {
	case conf.Bootstrap_UAT, conf.Bootstrap_Online:
		return 3
	default:
		return 4
	}
}

func GetGameSnakeID() int32 {
	switch conf.Env() {
	case conf.Bootstrap_UAT, conf.Bootstrap_Online:
		return 4
	default:
		return 6
	}
}

//口袋AR
func GetGameARID() int32 {
	switch conf.Env() {
	case conf.Bootstrap_UAT, conf.Bootstrap_Online:
		return 5
	default:
		return 7
	}
}

//神奇板子
func GetGameBoardID() int32 {
	return int32(10)
}

//水果合合合H5
func GetGameSGHHHIDH5() int32 {
	return 40
}

//中秋AR
func GetGameMidAutumnID() int32 {
	return 41
}

//年会AR
func GetGamePocketID() int32 {
	return 50
}

//拼图
func GetGamePTID() int32 {
	return 51
}

//宁波银行蛋黄小镇H5
func GetNBCBGameH5ID() int32 {
	return 52
}

//宁波银行蛋黄小镇小程序
func GetNBCBGameAppID() int32 {
	return 53
}

//中国银行H5（未上线）
func GetGameBocH5ID() int32 {
	return 54
}

//AR叠叠乐
func GetGameARDLL() int32 {
	return 60
}

func IsDigirlGame(gid uint64) bool {
	return gid == uint64(GetGameDfgirlID())
}

func IsBoardGame(gid uint64) bool {
	return gid == uint64(GetGameBoardID())
}
