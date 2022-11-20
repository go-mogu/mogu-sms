package util

import (
	"github.com/go-mogu/mgu-sms/global"
	"gopkg.in/gomail.v2"
)

func SendMail(subject, receiver, text string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(global.Cfg.Mail.UserName, global.Cfg.Mail.UserName))
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	d := gomail.NewDialer(global.Cfg.Mail.Host, global.Cfg.Mail.Port, global.Cfg.Mail.UserName, global.Cfg.Mail.Password)
	err := d.DialAndSend(m)
	return err
}
