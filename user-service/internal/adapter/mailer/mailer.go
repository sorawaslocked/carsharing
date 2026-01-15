package mailer

import (
	"context"
	"fmt"
	"github.com/mailersend/mailersend-go"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/mailer"
	"time"
)

type Mailer struct {
	ms   *mailersend.Mailersend
	from mailersend.From
}

func New(cfg mailer.Config) *Mailer {
	ms := mailersend.NewMailersend(cfg.APIKey)
	from := mailersend.From{
		Name:  "Car Rental",
		Email: cfg.From,
	}

	return &Mailer{
		ms:   ms,
		from: from,
	}
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
	message.SetFrom(m.from)
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
