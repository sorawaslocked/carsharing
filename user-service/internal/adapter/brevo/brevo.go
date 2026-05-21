package brevo

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	"carsharing/user-service/internal/model"
	brevocfg "carsharing/user-service/internal/pkg/brevo"
)

const sendURL = "https://api.brevo.com/v3/smtp/email"

type Brevo struct {
	log    *slog.Logger
	client *http.Client
	apiKey string
	from   emailAddr
}

type emailAddr struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type sendRequest struct {
	Sender      emailAddr   `json:"sender"`
	To          []emailAddr `json:"to"`
	Subject     string      `json:"subject"`
	TextContent string      `json:"textContent"`
	HTMLContent string      `json:"htmlContent"`
}

func New(log *slog.Logger, cfg brevocfg.Config) *Brevo {
	return &Brevo{
		log:    pkglog.WithComponent(log, "adapter.Brevo"),
		client: &http.Client{},
		apiKey: cfg.APIKey,
		from:   emailAddr{Email: cfg.From, Name: cfg.FromName},
	}
}

func (m *Brevo) SendActivationCode(ctx context.Context, receiver, code string) error {
	subject := "Your Activation Code"
	text := activationCodeText(code)
	html := activationCodeHTML(code)

	return m.send(ctx, receiver, subject, text, html)
}

func (m *Brevo) send(ctx context.Context, to, subject, text, html string) error {
	log := pkglog.WithMetadata(pkglog.WithMethod(m.log, "send"), utils.MetadataFromCtx(ctx))

	payload := sendRequest{
		Sender:      m.from,
		To:          []emailAddr{{Email: to}},
		Subject:     subject,
		TextContent: text,
		HTMLContent: html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Error("marshalling request", pkglog.Err(err))
		return model.ErrMailer
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sendURL, bytes.NewReader(body))
	if err != nil {
		log.Error("building request", pkglog.Err(err))
		return model.ErrMailer
	}
	req.Header.Set("api-key", m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		log.Error("sending request", pkglog.Err(err))
		return model.ErrMailer
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Error("unexpected response", slog.String("to", to), slog.Int("status", resp.StatusCode))
		return model.ErrMailer
	}

	return nil
}
