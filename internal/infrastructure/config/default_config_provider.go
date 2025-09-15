package config

type DefaultProvider struct{}

func NewDefaultProvider() *DefaultProvider {
	return &DefaultProvider{}
}

func (d *DefaultProvider) setValues(c *Config) error {
	c.Host = "localhost:8080"
	c.DatabaseDSN = ""
	c.AccrualSystemAddress = "localhost:8080"
	c.AccessTokenSecret = "my_secret_key"
	return nil
}
