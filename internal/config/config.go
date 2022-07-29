package config

import (
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"io/fs"
)

type App struct {
	LogFile            string `env:"LOG_FILE" envDefault:"D:\\Export\\Log\\killer1c77.log"`
	ProcessNamePattern string `env:"PROCESS_NAME_PATTERN" envDefault:"^1cv7"`
	HTTPPort           string `env:"HTTP_PORT" envDefault:"8081"`
}

func New() (cfg App, err error) {
	err = godotenv.Load()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return
	}

	err = env.Parse(&cfg)

	return
}
