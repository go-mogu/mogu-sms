package global

import (
	"github.com/go-mogu/hz-framework/config"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	Redis *redis.Client // redis连接池
	Cfg   *config.Conf  // yaml配置
	Viper *viper.Viper
)
