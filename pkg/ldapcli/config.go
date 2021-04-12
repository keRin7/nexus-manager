package ldapcli

type Config struct {
	BaseDN     string `env:"BASE_DN"`
	LdapServer string `env:"LDAP_SERVER"`
	//	BindDN       string `env:"BIND_DN"`
	//	Port         string `env:"PORT"`
	//	Host         string `env:"HOST"`
	//	BindPassword string `env:"BIND_PASSWORD"`
	//	Filter       string `env:"FILTER"`
}

func NewConfig() *Config {
	return &Config{}
}
