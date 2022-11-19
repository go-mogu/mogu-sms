package config

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-mogu/hz-framework/pkg/util"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"strings"
)

// ConfEnv env环境变量
var ConfEnv string

type (
	Conf struct {
		Server Server          `yaml:"server"`
		Nacos  NacosProperties `yaml:"nacos"`
		Zap    util.Zap        `yaml:"zap"`
		Redis  Redis           `yaml:"redis"`
		Amqp   Amqp            `yaml:"amqp"`
		Mail   Mail            `yaml:"mail"`
	}
	NacosProperties struct {
		Client    constant.ClientConfig
		Server    []constant.ServerConfig
		Config    vo.ConfigParam
		Discovery vo.RegisterInstanceParam
	}
	Server struct {
		Name string `yaml:"name"`
		Port string `yaml:"port"`
	}
	Redis struct {
		Host        string `yaml:"host" default:"127.0.0.1"`
		Port        string `yaml:"port" default:"6379"`
		Password    string `yaml:"password"`
		DbNum       int    `yaml:"dbNum" default:"0"`
		LoginPrefix string `yaml:"loginPrefix" default:"login_auth_"`
	}
	Amqp struct {
		Host     string `yaml:"host" default:"127.0.0.1"`
		Port     string `yaml:"port" default:"5672"`
		User     string `yaml:"user" default:"guest"`
		Password string `yaml:"password" default:""`
		Vhost    string `yaml:"vhost" default:""`
	}
	Mail struct {
		UserName string `yaml:"userName" default:"mogublog@163.com"`
		Password string `yaml:"password" default:"a1313375"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Auth     bool   `yaml:"auth"`
	}
)

//go:embed yaml
var yamlCfg embed.FS

// InitConfig 初始化配置
func InitConfig() (*Conf, *viper.Viper, error) {
	var cfg *Conf
	v := viper.New()
	v.SetConfigType("yaml")
	yamlConf, _ := yamlCfg.ReadFile("yaml/config." + ConfEnv + ".yaml")
	if err := v.ReadConfig(bytes.NewBuffer(yamlConf)); err != nil {
		return nil, v, err
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, v, err
	}
	cfg, err := InitNacosConfig(cfg, v)
	if err != nil {
		return nil, v, err
	}

	return cfg, v, nil
}

func InitNacosConfig(cfg *Conf, v *viper.Viper) (*Conf, error) {
	// 创建动态配置客户端
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cfg.Nacos.Client,
			ServerConfigs: cfg.Nacos.Server,
		},
	)
	if err != nil {
		return nil, err
	}
	cfg.Nacos.Config = getConfigParam(cfg.Server, cfg.Nacos.Config)

	content, err := configClient.GetConfig(cfg.Nacos.Config)
	err = v.ReadConfig(bytes.NewBuffer([]byte(content)))
	if err != nil {
		return nil, err
	}
	if err = v.Unmarshal(&cfg); err != nil {
		hlog.Error(err)
		return nil, err
	}
	return cfg, err
}

func getConfigParam(app Server, config vo.ConfigParam) vo.ConfigParam {
	config.DataId = fmt.Sprintf("%s-%s.%s", app.Name, config.Group, strings.ToLower(config.Type))
	config.Type = strings.ToUpper(config.Type)
	return config
}
