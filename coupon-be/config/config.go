package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Log   `yaml:"logger"`
		MYSQL `yaml:"mysql"`
		Cors  `yaml:"cors"`
		Redis `yaml:"redis"`
	}

	// App -.
	App struct {
		Name    string `yaml:"name"    env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// MYSQL -.
	MYSQL struct {
		PoolMax int    `yaml:"pool_max" env:"MYSQL_POOL_MAX"`
		URL     string `env-required:"false" env:"MYSQL_URL"`
	}

	// Cors -.
	Cors struct {
		AllowedOrigins []string `yaml:"allowed_origins"`
	}

	// Redis -.
	Redis struct {
		RedisAddress  string `yaml:"redis_address" env:"REDIS_ADDRESS"`
		RedisPassword string `yaml:"redis_password" env:"REDIS_PASSWORD"`
		RedisDB       int    `yaml:"redis_db" env:"REDIS_DB"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	_, er := os.ReadFile(".env")
	if er != nil {
		panic(er)
	}

	err := cleanenv.ReadConfig(".env", cfg)

	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
