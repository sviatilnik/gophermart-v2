package config

type Config struct {
	Host                 string
	DatabaseDSN          string
	AccrualSystemAddress string
	AccessTokenSecret    string
}

func NewConfig(providers ...Provider) Config {
	conf := Config{}
	for _, provider := range providers {
		_ = conf.setValues(provider)
	}

	return conf
}

func (c *Config) setValues(provider Provider) error {
	return provider.setValues(c)
}
