package listener

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-mogu/mgu-sms/global"
	"github.com/go-mogu/mgu-sms/internal/amqp/consumer"
	"github.com/go-mogu/mgu-sms/internal/consts"
	"github.com/go-mogu/mgu-sms/internal/handler"
	"github.com/go-mogu/mgu-sms/pkg/mq"
	"time"
)

// BlogLinter blog监听器
func BlogLinter() {
	amqpConfig := &mq.Config{
		User:     global.Cfg.Amqp.User,
		Password: global.Cfg.Amqp.Password,
		Host:     global.Cfg.Amqp.Host,
		Port:     global.Cfg.Amqp.Port,
		Vhost:    global.Cfg.Amqp.Vhost,
	}
	// 实例化amqp
	amqp := mq.New(amqpConfig, consts.MoguBlog, consts.ExchangeDirect, consts.MoguBlog, 1, 1, true)
	time.Sleep(time.Second * 1)
	if err := consumer.New(amqp, nil, handler.BlogHandler.BlogConsumerHandler).Consumer(); err != nil {
		hlog.Error("程序载入Mq服务失败！", err)
	}
}
