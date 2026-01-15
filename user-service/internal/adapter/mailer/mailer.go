package mailer

import (
	"context"
	"fmt"
	"github.com/mailersend/mailersend-go"
	"time"
)

var (
	from = mailersend.From{
		Name:  "Car Rental",
		Email: "car-rental@mail.com",
	}
)

type Mailer struct {
	ms *mailersend.Mailersend
}

func New(apiKey string) *Mailer {
	ms := mailersend.NewMailersend(apiKey)

	return &Mailer{ms: ms}
}

func (m *Mailer) SendActivationCode(ctx context.Context, to, code string) error {
	c, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	subject := "Activation Code"
	text := fmt.Sprintf("Your activation code: %s", code)
	html := fmt.Sprintf("Your activation code: %s", code)

	recipients := []mailersend.Recipient{
		{
			Email: to,
		},
	}

	message := m.ms.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)

	_, err := m.ms.Email.Send(c, message)
	if err != nil {
		return err
	}

	return nil
}
