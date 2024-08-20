package mail

import (
	"crypto/tls"
	"fmt"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/try"
	"github.com/gflydev/core/utils"
	"net/smtp"
)

type Envelop struct {
	To      []string // Required
	ReplyTo []string
	Bcc     []string
	Cc      []string
	Subject string // Required
	Text    string // Required
	HTML    string // Required
}

type protocol string

var smtpProtocol = protocol("smtp")

func Send(envelop Envelop) {
	protocol := utils.Getenv("MAIL_PROTOCOL", "smtp")

	if protocol != string(smtpProtocol) {
		log.Panicf("No support mail protocol `%s`", protocol)
	}

	e := New()
	e.From = fmt.Sprintf("%s <%s>",
		utils.Getenv("MAIL_NAME", "gFly - No Reply"),
		utils.Getenv("MAIL_SENDER", "no-reply@gfly.dev"),
	)

	if len(envelop.ReplyTo) == 0 {
		e.ReplyTo = []string{utils.Getenv("MAIL_SENDER", "no-reply@gfly.dev")}
	} else {
		e.ReplyTo = envelop.ReplyTo
	}

	if len(envelop.Bcc) > 0 {
		e.Bcc = envelop.Bcc
	}

	if len(envelop.Cc) > 0 {
		e.Cc = envelop.Cc
	}

	e.To = envelop.To
	e.Subject = envelop.Subject
	e.Text = []byte(envelop.Text)
	e.HTML = []byte(envelop.HTML)

	host := utils.Getenv("MAIL_HOST", "localhost")
	address := fmt.Sprintf("%s:%d", host, utils.Getenv("MAIL_PORT", 587))
	username := utils.Getenv("MAIL_USERNAME", "")
	password := utils.Getenv("MAIL_PASSWORD", "")

	try.Perform(func() {
		var err error
		auth := smtp.PlainAuth("", username, password, host)

		isTLS := utils.Getenv("MAIL_TLS", true)
		if isTLS {
			// TLS config
			tlsConfig := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         host,
			}
			err = e.SendWithStartTLS(address, auth, tlsConfig)
		} else {
			err = e.Send(address, auth)
		}

		if err != nil {
			try.Throw(err)
		}
	}).Catch(func(e try.E) {
		log.Errorf("Error send mail %v", e)
	})
}
