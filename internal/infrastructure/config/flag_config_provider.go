package config

import (
	"flag"
	"strings"
)

type FlagProvider struct {
	HostFlagName                 string
	DatabaseDSNFlagName          string
	AccrualSystemAddressFlagName string
}

func NewFlagProvider() *FlagProvider {
	return &FlagProvider{
		HostFlagName:                 "a",
		DatabaseDSNFlagName:          "d",
		AccrualSystemAddressFlagName: "r",
	}
}

func (flagConf *FlagProvider) setValues(c *Config) error {
	host := flag.String(flagConf.HostFlagName, "", "Адрес и порт запуска сервиса")
	databaseDSN := flag.String(flagConf.DatabaseDSNFlagName, "", "Адрес подключения к базе данных")
	accrualSystemAddress := flag.String(flagConf.AccrualSystemAddressFlagName, "", "Адрес системы расчёта начислений")
	flag.Parse()

	if strings.TrimSpace(*host) != "" {
		c.Host = *host
	}

	if strings.TrimSpace(*databaseDSN) != "" {
		c.DatabaseDSN = *databaseDSN
	}

	if strings.TrimSpace(*accrualSystemAddress) != "" {
		c.AccrualSystemAddress = *accrualSystemAddress
	}

	return nil
}
