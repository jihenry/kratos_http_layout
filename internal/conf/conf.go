package conf

import (
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/encoding"
	"gitlab.yeahka.com/gaas/pkg/util"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

var (
	global   config.Config
	boot     *Bootstrap
	confPath string
)

func Load(path string) (*Bootstrap, error) {
	global = config.New(
		config.WithSource(
			file.NewSource(path),
		),
		config.WithDecoder(func(src *config.KeyValue, target map[string]interface{}) error {
			if codec := encoding.GetCodec(src.Format); codec != nil {
				return codec.Unmarshal(src.Value, &target)
			}
			return nil
		}),
	)
	confPath = path
	if err := global.Load(); err != nil {
		return nil, err
	}
	var bc Bootstrap
	if err := global.Scan(&bc); err != nil {
		panic(err)
	}
	boot = &bc
	return boot, nil
}

//配置目录，xlsx使用
func GetConfDir() string {
	path, err := filepath.Abs(confPath)
	if err != nil {
		return ""
	}
	fi, err := os.Stat(path)
	if err != nil {
		return ""
	}
	if fi.IsDir() {
		return path
	}
	return filepath.Dir(path)
}

func Config() config.Config {
	return global
}

func Close() {
	if !util.IsNil(global) {
		global.Close()
	}
}

func TACfg() *TA {
	return boot.Ta
}

func Env() Bootstrap_Env {
	return boot.Env
}

func CacheCfg() *Data_Cache {
	if boot.Data.Cache != nil {
		return boot.Data.Cache
	}
	return &Data_Cache{}
}

func CronCfg() *Cron {
	if boot.Cron != nil {
		return boot.Cron
	}
	return &Cron{}
}

func LockCfg() *Lock {
	if boot.Lock != nil {
		return boot.Lock
	}
	return &Lock{}
}

func MiddlewareCfg() *Middleware {
	if boot.Middleware == nil {
		return &Middleware{}
	}
	return boot.Middleware
}

func NacosCfg() *Nacos {
	if boot.Nacos == nil {
		return &Nacos{}
	}
	return boot.Nacos
}

func KafkaCfg() *Kafka {
	if boot.Kafka == nil {
		return &Kafka{}
	}
	return boot.Kafka
}

func KafkaConsumerCfg() *KafkaConsumer {
	kafkaCfg := KafkaCfg()
	if kafkaCfg.Consumer == nil {
		return &KafkaConsumer{}
	}
	return kafkaCfg.Consumer
}

func KafkaProducerCfg() *KafkaProducer {
	kafkaCfg := KafkaCfg()
	if kafkaCfg.Producer == nil {
		return &KafkaProducer{}
	}
	return kafkaCfg.Producer
}

func TraceCfg() *Trace {
	if boot.Trace == nil {
		return &Trace{}
	}
	return boot.Trace
}

func MonitorCfg() *Monitor {
	if boot.Monitor == nil {
		return &Monitor{}
	}
	return boot.Monitor
}

func UserServiceCfg() *UserService {
	serviceCfg := ServiceCfg()
	if serviceCfg.User == nil {
		return &UserService{}
	}
	return serviceCfg.User
}

func ServiceCfg() *Service {
	if boot.Service == nil {
		return &Service{}
	}
	return boot.Service
}

func BaiduCfg() *BaiduMap {
	tpCfg := ThirdParyCfg()
	if tpCfg.Baidu == nil {
		return &BaiduMap{}
	}
	return tpCfg.Baidu
}

func RobotCfg() *WeixinRobot {
	tpCfg := ThirdParyCfg()
	if tpCfg.Robot == nil {
		return &WeixinRobot{}
	}
	return tpCfg.Robot
}

func ThirdParyCfg() *ThirdPary {
	if boot.Tp == nil {
		return &ThirdPary{}
	}
	return boot.Tp
}
