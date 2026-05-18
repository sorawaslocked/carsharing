package brevo

type Config struct {
	APIKey   string `yaml:"api_key"   env:"BREVO_API_KEY"   env-required:"true"`
	From     string `yaml:"from"      env:"BREVO_FROM"      env-required:"true"`
	FromName string `yaml:"from_name" env:"BREVO_FROM_NAME" env-default:"CarSharing"`
}
