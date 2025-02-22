package smtputils

import (
	"context"
	"sync"

	"github.com/jordanmarcelino/learn-go-microservices/pkg/config"
	"gopkg.in/gomail.v2"
)

var (
	dialer   *gomail.Dialer
	initOnce sync.Once
)

type Mailer interface {
	SendMail(ctx context.Context, to, subject, body string) error
}

type mailerImpl struct{}

func NewMailer() Mailer {
	initOnce.Do(func() {
		dialer = gomail.NewDialer(config.SMTP_CONFIG.Host, config.SMTP_CONFIG.Port, "", "")
	})

	return &mailerImpl{}
}

func (s *mailerImpl) SendMail(ctx context.Context, to, subject, body string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		m := gomail.NewMessage()

		m.SetHeader("From", config.SMTP_CONFIG.Email)
		m.SetHeader("To", to)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)

		return dialer.DialAndSend(m)
	}
}
