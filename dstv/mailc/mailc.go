package mailc

import (
	"log"

	"gopkg.in/gomail.v2"
)

type SMTP struct {
	Host     string
	Port     int
	Password string
	Username string
	logger   log.Logger
}

func NewSMTP(logger log.Logger, host, password, username string, port int) SMTP {
	smtp := SMTP{host, port, password, username, logger}
	return smtp
}
func (s SMTP) Send(subject, body, file string, to []string) error {
	mailer := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)
	ses, err := mailer.Dial()
	if err != nil {
		panic(err)
	}
	message := gomail.NewMessage()
	for _, e := range to {
		message.SetHeader("From", "no-reply@trinitytechnology.com")
		message.SetAddressHeader("To", e, e)
		message.SetHeader("Subject", "New Customer Payment from Mgurush")
		message.SetBody("text/plain", body)

		if file != "" {
			message.Attach(file)
		}

		if err := gomail.Send(ses, message); err != nil {
			s.logger.Println(err)
			return err
		}
	}
	message.Reset()
	return nil
}
