package mailer

type Config struct {
	APIKey string `env:"MAILER_SEND_API_KEY" env-required:"true"`
	From   string `env:"MAILER_SEND_FROM" env-required:"true"`
}
