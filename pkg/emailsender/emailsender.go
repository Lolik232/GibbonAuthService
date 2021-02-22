package emailsender

import (
	errors "auth-server/pkg/errors/types"
	"context"
	"fmt"
	"net/smtp"
	"net/textproto"
)

type IEmailSender interface {
	Send(ctx context.Context, subject, email, msgtype, msg string) error
}

type EmailSender struct {
	auth                                  smtp.Auth
	host, port, companyName, companyEmail string
}

func New(email, password, host, port, companyName, companyEmail string) IEmailSender {
	auth := smtp.PlainAuth("", email, password, host)
	return &EmailSender{
		auth,
		host,
		port,
		companyName,
		companyEmail,
	}
}
func (e EmailSender) Send(ctx context.Context, subject, email, msgtype, msg string) error {

	if len(subject) == 0 || len(email) == 0 || len(msgtype) == 0 || len(msg) == 0 {
		return errors.ErrInvalidArgument.New("Err. Params not be null.")
	}
	headers := make(map[string]string)
	headers["From"] = e.companyEmail
	headers["To"] = email
	headers["Subject"] = subject
	headers["Content-Type"] = msgtype
	message := ""
	for k, h := range headers {
		message += fmt.Sprintf("%s: %s\n", k, h)
	}
	message += "\n\r" + msg
	if err := smtp.SendMail(e.host+e.port, e.auth, e.companyEmail, []string{email}, []byte(message)); err != nil {
		if _, ok := err.(*textproto.Error); ok {
			return errors.ErrInvalidArgument.Newf("Invalid email %s", email)
		}
		return errors.NoType.New("")
	}
	return nil
}
