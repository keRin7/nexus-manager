package nexusmanager

import (
	"github.com/keRin7/nexus-manager/pkg/ldapcli"
)

type Config struct {
	Nexus_url      string `env:"NEXUS_URL"`
	Nexus_username string `env:"NEXUS_USERNAME"`
	Nexus_password string `env:"NEXUS_PASSWORD"`
	Nexus_repo     string `env:"NEXUS_REPO"`
	Ldap           *ldapcli.Config
}

func NewConfig() *Config {
	return &Config{
		Ldap: ldapcli.NewConfig(),
	}
}
