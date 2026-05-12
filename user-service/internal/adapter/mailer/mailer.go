package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-user-service/internal/pkg/log"
	mailercfg "github.com/sorawaslocked/car-rental-user-service/internal/pkg/mailer"
)

const sendURL = "https://send.api.mailtrap.io/api/send"

type Mailer struct {
	log    *slog.Logger
	client *http.Client
	token  string
	from   emailAddr
}

type emailAddr struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type sendRequest struct {
	From    emailAddr   `json:"from"`
	To      []emailAddr `json:"to"`
	Subject string      `json:"subject"`
	Text    string      `json:"text"`
	HTML    string      `json:"html"`
}

func New(log *slog.Logger, cfg mailercfg.Config) *Mailer {
	return &Mailer{
		log:    pkglog.WithComponent(log, "adapter.Mailer"),
		client: &http.Client{},
		token:  cfg.Token,
		from:   emailAddr{Email: cfg.From, Name: cfg.FromName},
	}
}

func (m *Mailer) SendActivationCode(ctx context.Context, receiver, code string) error {
	subject := "Your Activation Code"
	text := activationCodeText(code)
	html := activationCodeHTML(code)
	return m.send(ctx, receiver, subject, text, html)
}

func (m *Mailer) send(ctx context.Context, to, subject, text, html string) error {
	payload := sendRequest{
		From:    m.from,
		To:      []emailAddr{{Email: to}},
		Subject: subject,
		Text:    text,
		HTML:    html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		m.log.Error("marshalling request", pkglog.Err(err))
		return model.ErrMailer
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sendURL, bytes.NewReader(body))
	if err != nil {
		m.log.Error("building request", pkglog.Err(err))
		return model.ErrMailer
	}
	req.Header.Set("Authorization", "Bearer "+m.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.log.Error("sending request", pkglog.Err(err))
		return model.ErrMailer
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		m.log.Error("unexpected response", slog.String("to", to), slog.Int("status", resp.StatusCode))
		return model.ErrMailer
	}

	return nil
}

func activationCodeText(code string) string {
	return fmt.Sprintf(
		"Car Rental — Email Verification\n\n"+
			"Your activation code is:\n\n"+
			"  %s\n\n"+
			"Enter this code to verify your email address. It expires shortly.\n\n"+
			"If you did not request this, you can safely ignore this message.",
		code,
	)
}

func activationCodeHTML(code string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background:#f4f4f5;font-family:Arial,sans-serif">
  <table width="100%%" cellpadding="0" cellspacing="0" style="padding:40px 0">
    <tr><td align="center">
      <table width="480" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,.08)">
        <tr><td style="background:#111827;padding:24px 32px">
          <p style="margin:0;color:#ffffff;font-size:18px;font-weight:bold">Car Rental</p>
        </td></tr>
        <tr><td style="padding:32px">
          <p style="margin:0 0 8px;font-size:22px;font-weight:bold;color:#111827">Verify your email</p>
          <p style="margin:0 0 24px;font-size:14px;color:#6b7280">Use the code below to complete your registration.</p>
          <div style="background:#f9fafb;border:1px solid #e5e7eb;border-radius:6px;padding:20px;text-align:center;margin-bottom:24px">
            <span style="font-size:32px;font-weight:bold;letter-spacing:8px;color:#111827">%s</span>
          </div>
          <p style="margin:0;font-size:13px;color:#9ca3af">This code expires shortly. If you did not request this, you can safely ignore this email.</p>
        </td></tr>
      </table>
    </td></tr>
  </table>
</body>
</html>`, code)
}
