package errno

//下面这几个错误码是老错误码，没有做转换，为了这些错误码中间件和登录成功前有用到
var (
	ErrWrongParam   = new("ERR_WRONG_PARAM", 40000, "参数错误")   // 参数错误
	ErrUnauthorized = new("ERR_UNAUTHORIZED", 40001, "未登录")   // 无权限/未登录
	ErrInvalidSign  = new("ERR_INVALID_SIGN", 40002, "非法签名")  // 非法签名
	ErrReqRepeat    = new("ERR_REQ_REPEAT", 40010, "重复请求")    // 重复请求
	ErrExternal     = new("ERR_EXTERNAL", 40005, "下游请求异常")    // 下游请求异常
	ErrUnknown      = new("ERR_UNKNOWN", 40006, "未知错误")       // 未知错误
	ErrSystemError  = new("ERR_SYSTEM_ERROR", 50000, "服务器错误") // 服务器错误

	ErrEndActivity       = new("END_ACTIVITY_ERROR", 60000, "活动已结束")
	ErrDrawNum           = new("DRAW_NUM_ERROR", 60001, "抽奖次数不满足")
	ErrPrizeServiceError = new("PRIZE_SERVICE_ERROR", 65000, "奖品发放失败")
)

//80000 - 80099 游戏错误码
var (
	ErrNoDataFound            = newOldError("ERR_NO_DATA_FOUND", 40004, 80023, "数据不存在")               // 数据不存在
	ErrWrongPassWordOrAccount = newOldError("ERR_WRONG_PASSWORD_OR_ACCOUNT", 40007, 80026, "账号或密码错误") // 账号密码错误
	ErrWrongNoSuchLevel       = newOldError("ERR_WRONG_NO_SUCH_LEVEL", 40008, 80027, "没有这个等级")        // 账号密码错误
	ErrDataWrite              = newOldError("ERR_DATA_WRITE", 40009, 80028, "写入数据错误")                 // 写入数据错误
	ErrConfigInvalid          = newOldError("ERR_CONFIG_INVALID", 40011, 80030, "配置错误")               // 配置错误
	ErrItemNoEnough           = newOldError("ERR_ITEM_NOENOUGH", 40012, 80031, "道具不足")                // 道具不足
	ErrTopicLimitError        = newOldError("ERR_TOPIC_LIMIT", 40013, 80032, "今日答题次数已用尽")
	ErrRunIdInvalid           = newOldError("ERR_RUNID_INVALID", 40014, 80033, "runid非法")
	ErrTeamNumberLimitError   = newOldError("ERR_TEAM_NUMBER_LIMIT", 40015, 80034, "学号或工号已被使用")
	ErrConfigParamError       = newOldError("ERR_CONFIG_PARAM_ERROR", 50001, 80059, "配置参数错误")
	ErrInvalidJsCode          = newOldError("ERR_INVALID_JS_CODE", 40100, 80040, "jscode无效")           // jscode无效 ，过期
	ErrJsCodeRequestLimit     = newOldError("ERR_JS_CODE_REQUEST_LIMIT", 40101, 80041, "jscode请求过于频繁") // jscode请求过于频繁
	ErrWechatSystemBusy       = newOldError("ERR_WECHAT_SYSTEM_BUSY", 40102, 80042, "微信系统繁忙")          // 微信系统繁忙
	ErrBindWxGroupExist       = newOldError("ERR_BIND_WXGROUP_EXIST", 40103, 80043, "已经入群")            // 已经入群
	ErrBindWxGroupCode        = newOldError("ERR_BIND_WXGROUP_CODE", 40104, 80044, "邀请码不合法")           // 邀请码必须是在同一个服务器
	ErrWXGroupNoSign          = newOldError("ERR_WXGROUP_NO_SIGN", 40105, 80045, "还未群签到")              // 还未群签到
	ErrWXGroupSignTaked       = newOldError("ERR_WXGROUP_SIGN_TAKED", 40106, 80046, "群签到奖励已经领取")       // 群签到奖励已经领取
	ErrBindCustomerExist      = newOldError("ERR_BIND_CUSTOMER_EXIST", 40107, 80047, "已经是客服")          // 已经是客服
	ErrUserNotFound           = newOldError("ERR_USER_NOT_FOUND", 40108, 80048, "用户没有发现")              // 用户没有发现
	// 客户端
	ErrNoGameFound     = newOldError("ERR_NO_GAME_FOUND", 40200, 80050, "未找到游戏")      // 游戏错误
	ErrGameDataRange   = newOldError("ERR_GAME_DATA_RANGE", 40201, 80051, "数据超出范围")   // 修改游戏模块配置失败
	ErrShareBoxNoExist = newOldError("ERR_SBOX_NOEXIST", 40202, 80052, "没有宝箱")        // 没有宝箱
	ErrShareBoxTakeMax = newOldError("ERR_SBOX_TAKE_MAX", 40203, 80053, "今日领取宝箱已达上限") // 今日领取宝箱已达上限
	ErrShareBoxTaked   = newOldError("ERR_SBOX_TAKED", 40204, 80054, "已经领取此宝箱")       // 已经领取此宝箱
	ErrShareBoxSelf    = newOldError("ERR_SBOX_SELF", 40205, 80055, "不能领取自己的宝箱")      // 不能领取自己的宝箱
	ErrShareBoxUnknown = newOldError("ERR_SBOX_UNKNOWN", 40206, 80056, "未知错误")
	ErrDfgirlReviveMax = newOldError("ERR_DFGIRL_REVIVE_MAX", 40206, 80057, "跃动女孩复活超过最大次数") // 跃动女孩复活超过最大次数

	// 服务端
	//ServerError = newOldError("Server error", 50200) // 服务器错误
	// 种树
	ErrFriendTreeNotExists         = newOldError("ERR_FRIEND_TREE_NOT_EXISTS", 50501, 80060, "好友未开始种树") // 好友未种树
	ErrTreeModuleNotOpen           = newOldError("ERR_TREE_MODULE_NOT_OPEN", 50502, 80061, "种树入口未开启")   // 种树模块未开启
	ErrTreeHWaterMax               = newOldError("ERR_TREE_HWATER_MAX", 50503, 80062, "帮好友浇水达到最大值")     // 帮好友浇水达到最大值
	ErrTreeSeedPlanted             = newOldError("ERR_TREE_SEED_PLANTED", 50504, 80063, "种子已经种过")
	ErrTreePrizeInventoryNotEnough = newOldError("ERR_TREE_PRIZE_INVENTORY_NOT_ENOUGH", 50505, 80064, "种子对应奖励不足")
	ErrTreeActivityEnd             = newOldError("ERR_TREE_ACTIVITY_END", 50506, 80065, "种树活动已经结束")
	ErrTreeRoundBehind             = newOldError("ERR_TREE_ROUND_BEHIND", 50507, 80066, "种树轮次落后")

	// 贪吃蛇
	ErrSnakeShareTimeNotEnough = newOldError("ERR_SNAKE_SHARE_TIME_NOT_ENOUGH", 50601, 80070, "贪吃蛇分享次数不够") // 贪吃蛇分享次数不够
	ErrSnakeRepeatBuySkin      = newOldError("ERR_SNAKE_REPEAT_BUY_SKIN", 50602, 80071, "重复购买皮肤")
	ErrSnakeNoSuchSkin         = newOldError("ERR_SNAKE_NO_SUCH_SKIN", 50603, 80072, "没有这个皮肤")
	ErrSnakeNotEnoughAsset     = newOldError("ERR_SNAKE_NOT_ENOUGH_ASSET", 50604, 80073, "资产不足")

	// 消消乐
	ErrCrushLimitActivityNotExist           = newOldError("ERR_CRUSH_LIMIT_ACTIVITY_NOT_EXIST", 50701, 80080, "消消乐限时挑战活动不存在")
	ErrCrushLimitActivityShareTimeNotEnough = newOldError("ERR_CRUSH_LIMIT_ACTIVITY_SHARE_TIME_NOT_ENOUGH", 50702, 80081, "消消乐限时挑战分享次数不够")
	ErrCrushNoSuchCheckPoint                = newOldError("ERR_CRUSH_NO_SUCH_CHECKPOINT", 50703, 80082, "消消乐没有此关")
	ErrCrushLimitActivityPowerNotEnough     = newOldError("ERR_CRUSH_LIMIT_ACTIVITY_POWER_NOT_ENOUGH", 50704, 80083, "消消乐限时挑战体力不够")
	ErrCrushCantShareGet                    = newOldError("ERR_CRUSH_CANT_SHARE_GET", 50705, 80084, "当前不能同构分享获取")
	ErrCrushShareCountMax                   = newOldError("ERR_CRUSH_SHARE_COUNT_MAX", 50706, 80085, "分享获得以达到上限")

	// 口袋ar
	ErrPocketShareTimeNotEnough = newOldError("ERR_POCKET_SHARE_TIME_NOT_ENOUGH", 50801, 80090, "分享次数不够")
	ErrPocketPowerNotEnough     = newOldError("ERR_POCKET_POWER_NOT_ENOUGH", 50802, 80091, "体力不够")
	ErrPocketShareCodeInvalid   = newOldError("ERR_POCKET_SHARE_CODE_INVALID", 50803, 80092, "链接失效")
	ErrPocketShareCodeSelfClick = newOldError("ERR_POCKET_SHARE_CODE_SELF_CLICK", 50804, 80093, "不能自己点卡片")
)

// 80100 - 80199 业务错误码
var (
	ErrTmplUsedInCorpGameCfg           = newOldError("ERR_TMPL_USED_IN_CORP_GAME_CFG", 51210, 80100, "游戏模板已被定制，不可编辑")                          // 游戏模板已被定制，不可编辑
	ErrDuplicateNameInTmpl             = newOldError("ERR_DUPLICATE_NAME_IN_TMPL", 51211, 80101, "游戏模板名称重复")                                   // 游戏模板Name重复
	ErrGameModuleParam                 = newOldError("ERR_GAME_MODULE_PARAM", 51212, 80102, "游戏模块参数异常")                                        // 游戏模块参数异常
	ErrGameModuleInternal              = newOldError("ERR_GAME_MODULE_INTERNAL", 51213, 80103, "修改游戏模块配置失败")                                   // 修改游戏模块配置失败
	ErrPrizeInActivity                 = newOldError("ERR_PRIZE_IN_ACTIVITY", 51220, 80104, "奖品已被活动关联，不可编辑/删除")                                // 奖品已被活动关联，不可编辑/删除
	ErrBrandGameNameIsExists           = newOldError("ERR_BRAND_GAME_NAME_IS_EXISTS", 51230, 80105, "品牌游戏名称已存在（同一个品牌下）")                       // 品牌游戏名称已存在（同一个品牌下）
	ErrServerNameIsExists              = newOldError("ERR_SERVER_NAME_IS_EXISTS", 51240, 80106, "服务器名称已存在（同一个品牌游戏下）")                          // 服务器名称已存在（同一个品牌游戏下）
	ErrDuplicateActivityTypeInSameTime = newOldError("ERR_DUPLICATE_ACTIVITY_TYPE_IN_SAME_TIME", 51241, 80107, "相同类型的子活动在同一时段下不可重复创建（同一服务器下）") // 相同类型的子活动在同一时段下不可重复创建（同一服务器下）
	ErrRepeatBuySkin                   = newOldError("ERR_REPEAT_BUY_SKIN", 60101, 80110, "重复购买皮肤")
	ErrNotOwnSkin                      = newOldError("ERR_NOT_OWN_SKIN", 60102, 80111, "不拥有皮肤")
	ErrTeamNotExist                    = newOldError("ERR_TEAM_NOT_EXSIT", 60103, 80112, "团队不存在")
	ErrTeamJoinBefore                  = newOldError("ERR_TEAM_JOIN_BEFORE", 60104, 80113, "已经加入过团队")
	ErrTeamNotJoin                     = newOldError("ERR_TEAM_NOT_JOIN", 60105, 80114, "不在团队内")
	ErrPicOwnBefore                    = newOldError("ERR_PIC_OWN_BEFORE", 60106, 80115, "已经拥有图鉴")
	ErrBuySkinFail                     = newOldError("ERR_BUY_SKIN_FAIL", 60107, 80116, "购买不成功")
	ErrPluginOff                       = newOldError("ERR_PLUGIN_OFF", 80117, 80117, "插件已关闭")
	ErrGameSpaceItemError              = newOldError("ERR_GAMESPACE_ITEM_ERROR", 70000, 80118, "道具已经领取")
	ErrPingTuTimesNotEnough            = newOldError("今天赠送/索要的次数已达上限", 60203, 80119, "今天赠送/索要的次数已达上限")
	ErrPingTuShareSelf                 = newOldError("无法赠送拼图給自己", 60204, 80120, "无法赠送拼图給自己")
	ErrPingTuShareCardNotEnough        = newOldError("卡片不足", 60202, 80121, "卡片不足")
	ErrPingTuShareCardInvalid          = newOldError("分享卡片失效", 60201, 80122, "分享卡片失效")
)
