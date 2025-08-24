package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"server/global"
	"strings"
)

func Email(To, subject string, body string) error {
	to := strings.Split(To, ",")
	return send(to, subject, body)
}
func send(to []string, subject string, body string) error {
	emailCfg := global.Config.Email
	from := emailCfg.From
	nickname := emailCfg.Nickname
	secret := emailCfg.Secret
	host := emailCfg.Host
	port := emailCfg.Port
	isSSL := emailCfg.IsSSL

	auth := smtp.PlainAuth("", from, secret, host)

	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		e.From = from
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)

	var err error
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if isSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}

/*
Email 函数：（大写）
负责处理业务层面的逻辑（如参数校验、格式转换），例如将逗号分隔的邮箱地址转换为切片。
它是面向外部的接口，直接供其他代码调用。
send 函数：（小写）
专注于底层实现细节（如 SMTP 连接、邮件发送），不关心参数来源，只负责执行发送操作。
*/
