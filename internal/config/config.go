package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Port string `env:"BOOKING_PORT" env-default:"9449"`
	} `env-prefix:""`

	DB struct {
		DSN string `env:"DATABASE_DSN" env-required:"true"`
	} `env-prefix:""`

	Auth struct {
		JWTSecret string `env:"JWT_SECRET" env-required:"true"`
	} `env-prefix:""`

	Timeouts struct {
		DB time.Duration `env:"DB_TIMEOUT" env-default:"5s"`
	} `env-prefix:""`
}

func Load() Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}
	return cfg
}
