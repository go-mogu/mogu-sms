package bootstrap

import (
	"bytes"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-mogu/hz-framework/config"
	"github.com/go-mogu/hz-framework/global"
	"github.com/go-mogu/hz-framework/internal/listener"
	"github.com/go-mogu/hz-framework/pkg/lib"
	"github.com/go-mogu/hz-framework/pkg/log"
	"github.com/go-mogu/hz-framework/pkg/util"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// 定义服务列表
const (
	LoggerService = `Logger`
	RedisService  = `Redis`
	MqService     = `mq`
)

type bootServiceMap map[string]func() error

// BootedService 已经加载的服务
var (
	BootedService []string
	err           error
	// serviceMap 程序启动时需要自动加载的服务
	serviceMap = bootServiceMap{
		RedisService:  bootRedis,
		LoggerService: bootLogger,
		MqService:     bootMq,
	}
)

// BootService 加载服务
func BootService(services ...string) {
	// 初始化配置
	if err = bootConfig(); err != nil {
		panic("初始化config配置失败：" + err.Error())
	}
	if len(services) == 0 {
		services = serviceMap.keys()
	}
	BootedService = make([]string, 0)
	for k, val := range serviceMap {
		if util.InAnySlice[string](services, k) {
			if err := val(); err != nil {
				panic("程序服务启动失败:" + err.Error())
			}
			BootedService = append(BootedService, k)
		}
	}
}

// bootConfig 载入配置
func bootConfig() error {
	global.Cfg, global.Viper, err = config.InitConfig()
	if err == nil {
		err = ListenConfig()
	}
	return err
}

func ListenConfig() error {
	// 创建动态配置客户端
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &global.Cfg.Nacos.Client,
			ServerConfigs: global.Cfg.Nacos.Server,
		},
	)
	global.Cfg.Nacos.Config.OnChange = func(namespace, group, dataId, data string) {
		hlog.Debug("group:" + group + ", dataId:" + dataId + ", configure to change!")
		v := global.Viper
		err = v.ReadConfig(bytes.NewBuffer([]byte(data)))
		if err != nil {
			hlog.Error(err)
			return
		}
		if err := v.Unmarshal(&global.Cfg); err != nil {
			hlog.Error(err)
			return
		}
	}
	err = configClient.ListenConfig(global.Cfg.Nacos.Config)
	return err
}

// bootLogger 将配置载入日志服务
func bootLogger() error {
	logger := log.NewLogger(global.Cfg.Zap.Director)
	defer logger.Sync()
	hlog.SetLogger(logger)
	hlog.Infof("程序载入Logger服务成功 [ 日志路径：%s ]", global.Cfg.Zap.Director)
	return err
}

// bootRedis 装配redis服务
func bootRedis() error {
	redisConfig := lib.RedisConfig{
		Addr:     fmt.Sprintf("%s:%s", global.Cfg.Redis.Host, global.Cfg.Redis.Port),
		Password: global.Cfg.Redis.Password,
		DbNum:    global.Cfg.Redis.DbNum,
	}
	global.Redis, err = lib.NewRedis(redisConfig)
	if err == nil {
		hlog.Info("程序载入Redis服务成功")
	}
	return err
}

// bootMq 装配mq
func bootMq() error {
	listener.Init()
	return nil
}

// keys 获取BootServiceMap中所有键值
func (m bootServiceMap) keys() []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
