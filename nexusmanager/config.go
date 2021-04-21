package nexusmanager

import (
	"github.com/keRin7/nexus-manager/pkg/ldapclient"
)

type Config struct {
	Nexus_url      string   `env:"NEXUS_URL"`
	Nexus_username string   `env:"NEXUS_USERNAME"`
	Nexus_password string   `env:"NEXUS_PASSWORD"`
	Nexus_repo     string   `env:"NEXUS_REPO"`
	Admin_users    []string `env:"ADMIN_USERS" envSeparator:" "`
	Ldap           *ldapclient.Config
}

func NewConfig() *Config {
	return &Config{
		Ldap: ldapclient.NewConfig(),
	}
}
