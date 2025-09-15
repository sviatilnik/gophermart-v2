package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/gophermart/internal/infrastructure/config/mock_config"
	"go.uber.org/mock/gomock"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{
		os.Args[0],
		"-tta=https://google.com",
	}

	config := NewConfig(
		&DefaultProvider{},
		&FlagProvider{
			HostFlagName:                 "tta",
			DatabaseDSNFlagName:          "ttd",
			AccrualSystemAddressFlagName: "ttr",
		},
		NewEnvProvider(getMockEnvGetter(t)),
	)

	assert.Equal(t, "https://google.com", config.Host)             // from flag provider
	assert.Equal(t, "database_dsn", config.DatabaseDSN)            // from env provider
	assert.Equal(t, "localhost:8080", config.AccrualSystemAddress) // from default provider
}

func getMockEnvGetter(t *testing.T) EnvGetter {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_config.NewMockEnvGetter(ctrl)
	m.EXPECT().LookupEnv("DATABASE_URI").Return("database_dsn", true).AnyTimes()
	m.EXPECT().LookupEnv(gomock.Any()).Return("", false).AnyTimes()

	return m
}
