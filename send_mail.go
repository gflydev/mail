package mail

import (
	"crypto/tls"
	"errors"
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

// Send builds a message from the given Envelop using the MAIL_* environment
// variables and delivers it over SMTP. It returns any error encountered while
// sending so callers can react to delivery failures; the error is also logged
// for backward compatibility. Callers that don't care about the result may
// safely ignore the returned value.
func Send(envelop Envelop) error {
	protocol := utils.Getenv("MAIL_PROTOCOL", "smtp")

	if protocol != string(smtpProtocol) {
		err := fmt.Errorf("unsupported mail protocol %q", protocol)
		log.Errorf("Error send mail %v", err)
		return err
	}

	if len(envelop.To) == 0 {
		err := errors.New("envelop must specify at least one To address")
		log.Errorf("Error send mail %v", err)
		return err
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

	var sendErr error
	try.Perform(func() {
		var err error
		auth := smtp.PlainAuth("", username, password, host)

		isTLS := utils.Getenv("MAIL_TLS", true)
		if isTLS {
			// TLS config. Certificate verification can be disabled for
			// development servers (e.g. self-signed certs) via MAIL_TLS_SKIP_VERIFY.
			tlsConfig := &tls.Config{
				InsecureSkipVerify: utils.Getenv("MAIL_TLS_SKIP_VERIFY", false), //nolint:gosec // opt-in via env for dev servers
				ServerName:         host,
				MinVersion:         tls.VersionTLS12,
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
		sendErr = fmt.Errorf("%v", e)
	})

	return sendErr
}
