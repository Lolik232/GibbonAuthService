package emailsender

import (
	"context"
	"net/smtp"
)

type IEmailSender interface {
	Send(ctx context.Context, email, subject, msg string) error
}

type EmailSender struct {
	auth                            smtp.Auth
	port, companyName, companyEmail string
}

func New(email, password, host, port, companyName, companyEmail string) *EmailSender {
	auth := smtp.PlainAuth("", email, password, host)
	return &EmailSender{
		auth,
		host,
		companyName,
		companyEmail,
	}
}

func (e EmailSender) Send(ctx context.Context, email, msg string) error {

	panic("implement me")
}
