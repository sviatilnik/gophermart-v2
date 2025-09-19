package config

import (
	"os"
	"strings"
)

type EnvGetter interface {
	LookupEnv(key string) (string, bool)
}

type OSEnvGetter struct{}

func NewOSEnvGetter() *OSEnvGetter {
	return &OSEnvGetter{}
}

func (O *OSEnvGetter) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

type EnvProvider struct {
	getter EnvGetter
}

func NewEnvProvider(getter EnvGetter) *EnvProvider {
	return &EnvProvider{
		getter: getter,
	}
}

func (env *EnvProvider) setValues(c *Config) error {
	host, ok := env.getter.LookupEnv("RUN_ADDRESS")
	if ok && strings.TrimSpace(host) != "" {
		c.Host = host
	}

	databaseDSN, ok := env.getter.LookupEnv("DATABASE_URI")
	if ok && strings.TrimSpace(databaseDSN) != "" {
		c.DatabaseDSN = databaseDSN
	}

	accrualSystemAddress, ok := env.getter.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if ok && strings.TrimSpace(accrualSystemAddress) != "" {
		c.AccrualSystemAddress = accrualSystemAddress
	}

	return nil
}
