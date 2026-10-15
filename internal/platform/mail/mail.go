// Package mail provides SMTP delivery.
package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

// TODO: 读取配置，代码重构（发件人/收件人/SMTP 地址当前为占位符）
type MailDto struct {
	Subject string
	Body    string
}

func (m *MailDto) Send() (bool, error) {
	for {
		from := mail.Address{Name: "发件人名称", Address: "发件人邮箱地址"}
		to := mail.Address{Name: "收件人名称", Address: "收件人邮箱地址"}

		headers := make(map[string]string)
		headers["From"] = from.String()
		headers["To"] = to.String()
		headers["Subject"] = m.Subject

		message := ""
		for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + m.Body

		servername := "smtp.server:smtp.port"
		host, _, _ := net.SplitHostPort(servername)
		auth := smtp.PlainAuth("", "smtp.user", "smtp.password", host)

		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		conn, err := tls.Dial("tcp", servername, tlsconfig)
		if err != nil {
			break
		}

		c := &smtp.Client{}
		c, err = smtp.NewClient(conn, host)
		if err != nil {
			break
		}

		if err = c.Auth(auth); err != nil {
			break
		}

		if err = c.Mail(from.Address); err != nil {
			break
		}

		if err = c.Rcpt(to.Address); err != nil {
			break
		}

		w, err := c.Data()
		if err != nil {
			break
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			break
		}

		err = w.Close()
		if err != nil {
			break
		}

		c.Quit()

		break
	}

	return true, nil
}
