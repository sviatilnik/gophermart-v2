package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFlagProvider(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{
		os.Args[0],
		"-ta=https://google.com",
		"-td=database_dsn",
		"-tr=https://short.google.com",
	}
	config := NewConfig(&FlagProvider{
		HostFlagName:                 "ta",
		DatabaseDSNFlagName:          "td",
		AccrualSystemAddressFlagName: "tr",
	})

	assert.Equal(t, "https://google.com", config.Host)
	assert.Equal(t, "https://short.google.com", config.AccrualSystemAddress)
	assert.Equal(t, "database_dsn", config.DatabaseDSN)
}
