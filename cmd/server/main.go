package main

import (
	"flag"
	"hlayout/internal/conf"
	aesPkg "hlayout/internal/pkg/aes_pkg"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"time"

	"gitlab.yeahka.com/gaas/pkg/config"
	"gitlab.yeahka.com/gaas/pkg/govern"
	"gitlab.yeahka.com/gaas/pkg/mq/kafka"
	"gitlab.yeahka.com/gaas/pkg/registry"
	"gitlab.yeahka.com/gaas/pkg/rpc"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	plog "gitlab.yeahka.com/gaas/pkg/log"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "hlayout"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname() //TODO:自动获取
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(register *nacos.Registry, servers ...transport.Server) *kratos.App {
	if err := initPkg(); err != nil {
		panic(err)
	}
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Registrar(register),
		kratos.Server(servers...),
	)
}

func initLogger(c *conf.Zap) error {
	baseLogger, err := plog.NewZapLogger(
		plog.WithConsole(c.Console),
		plog.WithDir(c.Dir),
		plog.WithFileName(c.FileName),
		plog.WithLevel(c.Level),
	)
	if err != nil {
		return err
	}
	logger := log.With(baseLogger,
		"time", log.DefaultTimestamp,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
	log.SetLogger(logger) //防止直接调用karao的log打印
	plog.SetLogger(logger)
	return nil
}

func initNacos(c *conf.Nacos) (*nacos.Registry, error) {
	//1. 注册中心
	register, err := registry.NewNacosClient(c.Addr, uint64(c.Port),
		constant.WithCacheDir(c.CacheDir),
		constant.WithLogDir(c.LogDir),
		constant.WithNamespaceId(c.Namespace),
		constant.WithLogLevel(c.LogLevel),
		constant.WithNotLoadCacheAtStart(c.NotLoadCacheAtStart),
		constant.WithTimeoutMs(uint64(c.TimeoutMs)))
	if err != nil {
		return nil, err
	}
	rpc.SetDiscovery(register)
	registry.SetDiscovery(register)
	//2. 配置中心
	configClient, err := config.NewNacosClient(c.Addr, uint64(c.Port),
		constant.WithCacheDir(c.CacheDir),
		constant.WithLogDir(c.LogDir),
		constant.WithNamespaceId(c.Namespace),
		constant.WithLogLevel(c.LogLevel),
		constant.WithNotLoadCacheAtStart(c.NotLoadCacheAtStart),
		constant.WithTimeoutMs(uint64(c.TimeoutMs)))
	if err != nil {
		return nil, err
	}
	config.SetGlobalConfigClient(configClient)
	return register, nil
}

func initPkg() error {
	if err := aesPkg.Init(); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	//1. 加载本地配置
	bc, err := conf.Load(flagconf)
	defer conf.Close()
	if err != nil {
		panic(err)
	}

	//2. 初始化日志和基础库
	if err := initLogger(bc.Zap); err != nil {
		panic(err)
	}

	//3. 初始化nacos，创建远程配置终端，注册中心
	registry, err := initNacos(bc.Nacos)
	if err != nil {
		panic(err)
	}

	//5. 初始化链路追踪
	err = govern.InitTracerProviderUseAgent(bc.Trace.AgentHost, bc.Trace.AgentPort, bc.Trace.ServerName)
	if err != nil {
		panic(err)
	}

	//6. kafka生产者
	kafkaProducer, err := kafka.NewKafkaProducer(bc.Kafka.Addrs)
	if err != nil {
		panic(err)
	}
	kafka.SetGlobalProducer(kafkaProducer)

	//7. 初始化服务
	app, cleanup, err := wireApp(bc.Data, bc.Server, registry)
	if err != nil {
		panic(err)
	}

	defer cleanup()
	//8. 启动服务
	if err := app.Run(); err != nil {
		panic(err)
	}
}
