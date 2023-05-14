package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax        int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		VaTimezone     string `env-required:"true"   yaml:"pg_timezone"    env:"VA_TIMEZONE"`
		DatabaseDriver string `env-required:"true"   yaml:"database_driver"    env:"DATABASE_DRIVER"`
		PostgresUrl    string `env-required:"true"   yaml:"PG_URL"    env:"PG_URL"`
	}
)

// NewConfig returns app config -.
func NewConfig() (*Config, error) {

	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)

	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type T struct {
	Storage struct {
		File struct {
			Path string `json:"path"`
		} `json:"file"`
	} `json:"storage"`
	Listener struct {
		Tcp struct {
			Address    string `json:"address"`
			TlsDisable bool   `json:"tls_disable"`
		} `json:"tcp"`
	} `json:"listener"`
	Ui              bool   `json:"ui"`
	MaxLeaseTtl     string `json:"max_lease_ttl"`
	DefaultLeaseTtl string `json:"default_lease_ttl"`
	DisableMlock    bool   `json:"disable_mlock"`
}
