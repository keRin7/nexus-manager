package ldapcli

import (
	"crypto/tls"
	"log"

	"github.com/go-ldap/ldap/v3"
	"github.com/sirupsen/logrus"
)

type LdapCli struct {
	Config *Config
	Conn   *ldap.Conn
}

func New(config *Config) *LdapCli {
	return &LdapCli{
		Config: config,
	}
}

func (c *LdapCli) Init() {
	var err error

	c.Conn, err = ldap.DialURL(c.Config.LdapServer, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		log.Println(c.Config.LdapServer)
		log.Fatal(err)
	}
	//ldap.AddAttribute
	//cfg := &ldap.Config{
	//	BaseDN:     c.Config.BaseDN,
	//	LdapServer: c.Config.LdapServer,
	//}
	//c.authenticator = auth.New()
	//c.cache = store.NewFIFO(context.Background(), time.Minute*10)
	//strategy := ldap.NewCached(cfg, c.cache)
	//c.authenticator.EnableStrategy(ldap.StrategyKey, strategy)
	//c.authenticator.Authenticate()
}

func (c *LdapCli) TryToBind(username string, password string) bool {
	logrus.Println("cn=" + username + "," + c.Config.BaseDN)
	err := c.Conn.Bind("uid="+username+","+c.Config.BaseDN, password)
	if err != nil {
		logrus.Println("Bind error")
		logrus.Println(err)
		return false
	}
	logrus.Println("Bind success")
	return true
}
