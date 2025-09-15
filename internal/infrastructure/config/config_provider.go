package config

type Provider interface {
	setValues(c *Config) error
}
