package nexusmanager

type Config struct {
	Nexus_url      string `env:"NEXUS_URL"`
	Nexus_username string `env:"NEXUS_USERNAME"`
	Nexus_password string `env:"NEXUS_PASSWORD"`
	Nexus_repo     string `env:"NEXUS_REPO"`
}

func NewConfig() *Config {
	return &Config{}
}
