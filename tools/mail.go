package tools

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

// TODO: 读取配置，代码重构
type MailDto struct {
	Subject string
	Body    string
}

func (this *MailDto) Send() (bool, error) {
	for {
		from := mail.Address{"发件人名称", "发件人邮箱地址"}
		to := mail.Address{"收件人名称", "收件人邮箱地址"}

		// Setup headers
		headers := make(map[string]string)
		headers["From"] = from.String()
		headers["To"] = to.String()
		headers["Subject"] = this.Subject

		// Setup message
		message := ""
		for k, v := range headers {
			message += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		message += "\r\n" + this.Body

		// Connect to the SMTP Server
		servername := "smtp.server:smtp.port"
		host, _, _ := net.SplitHostPort(servername)
		auth := smtp.PlainAuth("", "smtp.user", "smtp.password", host)

		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		// Here is the key, you need to call tls.Dial instead of smtp.Dial
		// for smtp servers running on 465 that require an ssl connection
		// from the very beginning (no starttls)
		conn, err := tls.Dial("tcp", servername, tlsconfig)
		if err != nil {
			break
		}

		c := &smtp.Client{}
		c, err = smtp.NewClient(conn, host)
		if err != nil {
			break
		}

		// Auth
		if err = c.Auth(auth); err != nil {
			break
		}

		// To && From
		if err = c.Mail(from.Address); err != nil {
			break
		}

		if err = c.Rcpt(to.Address); err != nil {
			break
		}

		// Data
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
