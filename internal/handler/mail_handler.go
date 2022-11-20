package handler

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-mogu/mgu-sms/internal/util"
	"github.com/go-mogu/mgu-sms/pkg/mq"
	"github.com/go-mogu/mgu-sms/pkg/util/gconv"
)

type mailHandler struct{}

var MailHandler = &mailHandler{}

// MailConsumerHandler 处理消费者方法
func (s *mailHandler) MailConsumerHandler(mq *mq.RabbitMQ, data map[string]interface{}) error {
	chErrors := make(chan error)
	go func() {
		deliveries, err := mq.Consume()
		if err != nil {
			chErrors <- err
		}
		for d := range deliveries {
			hlog.Info("开始发送送邮件")
			// 消费mq，发送邮件
			data = gconv.Map(d.Body)
			err = util.SendMail(gconv.String(data["subject"]), gconv.String(data["receiver"]), gconv.String(data["text"]))
			if err != nil {
				break
			}
			err = d.Ack(true)
			if err != nil {
				return
			}
		}
	}()
	select {
	case err := <-chErrors:
		close(chErrors)
		hlog.Errorf("Consumer failed: %s\n", err)
		return err
	}
	return nil
}
