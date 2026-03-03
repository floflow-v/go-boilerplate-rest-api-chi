package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Api      ApiConfig      `envPrefix:"API_"`
	Log      LogConfig      `envPrefix:"LOG_"`
	Database DatabaseConfig `envPrefix:"DATABASE_"`
}

type ApiConfig struct {
	Environment string `env:"ENVIRONMENT,required,notEmpty"`
	Host        string `env:"HOST,required,notEmpty"`
	Port        int    `env:"PORT,required,notEmpty"`
}
type LogConfig struct {
	Level  string `env:"LEVEL,required,notEmpty"`
	Format string `env:"FORMAT,required,notEmpty"`
}

type DatabaseConfig struct {
	Host            string        `env:"HOST,required,notEmpty"`
	Port            int           `env:"PORT,required,notEmpty"`
	User            string        `env:"USER,required,notEmpty"`
	DBPassword      string        `env:"PASSWORD,required,notEmpty"`
	Name            string        `env:"NAME,required,notEmpty"`
	LogLevel        string        `env:"LOG_LEVEL,required,notEmpty"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS,required,notEmpty"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS,required,notEmpty"`
	MaxLifetimeConn time.Duration `env:"MAX_LIFETIME_CONN,required,notEmpty"`
	MaxIdleTimeConn time.Duration `env:"MAX_IDLE_TIME_CONN,required,notEmpty"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
