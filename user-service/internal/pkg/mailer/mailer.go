package mailer

type Config struct {
	Token    string `yaml:"token" env:"MAILTRAP_TOKEN" env-required:"true"`
	From     string `yaml:"from" env:"MAILTRAP_FROM" env-required:"true"`
	FromName string `yaml:"from_name" env:"MAILTRAP_FROM_NAME" env-default:"Car Rental"`
}
