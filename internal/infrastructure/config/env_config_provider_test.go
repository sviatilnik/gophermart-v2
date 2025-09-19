package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/sviatilnik/gophermart/internal/infrastructure/config/mock_config"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestEnvProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_config.NewMockEnvGetter(ctrl)

	m.EXPECT().LookupEnv("RUN_ADDRESS").Return("https://google.com", true).AnyTimes()
	m.EXPECT().LookupEnv("DATABASE_URI").Return("database_dsn", true).AnyTimes()
	m.EXPECT().LookupEnv("ACCRUAL_SYSTEM_ADDRESS").Return("https://google.com", true).AnyTimes()

	config := NewConfig(NewEnvProvider(m))

	assert.Equal(t, "https://google.com", config.Host)
	assert.Equal(t, "database_dsn", config.DatabaseDSN)
	assert.Equal(t, "https://google.com", config.AccrualSystemAddress)
}
