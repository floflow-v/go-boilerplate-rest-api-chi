package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go-rest-api-chi-example/internal/config"
)

func TestLoadConfig(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		t.Setenv("API_ENVIRONMENT", "test")
		t.Setenv("API_HOST", "localhost")
		t.Setenv("API_PORT", "8080")
		t.Setenv("LOG_LEVEL", "info")
		t.Setenv("LOG_FORMAT", "json")
		t.Setenv("DATABASE_HOST", "localhost")
		t.Setenv("DATABASE_PORT", "5432")
		t.Setenv("DATABASE_USER", "testuser")
		t.Setenv("DATABASE_PASSWORD", "testpass")
		t.Setenv("DATABASE_NAME", "testdb")
		t.Setenv("DATABASE_LOG_LEVEL", "warn")
		t.Setenv("DATABASE_MAX_OPEN_CONNS", "10")
		t.Setenv("DATABASE_MAX_IDLE_CONNS", "5")
		t.Setenv("DATABASE_MAX_LIFETIME_CONN", "1h")
		t.Setenv("DATABASE_MAX_IDLE_TIME_CONN", "10m")

		newCfg, err := config.LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, "test", newCfg.Api.Environment)
		assert.Equal(t, "localhost", newCfg.Api.Host)
		assert.Equal(t, 8080, newCfg.Api.Port)
	})

	t.Run("assert error", func(t *testing.T) {
		newCfg, err := config.LoadConfig()

		assert.Error(t, err)
		assert.Equal(t, config.Config{}, newCfg)
	})
}
