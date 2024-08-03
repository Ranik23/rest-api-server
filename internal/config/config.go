package config

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}


type HTTPServer struct {
	Address string `yaml:"address" env-default:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}


func MustLoad() *Config{

	configPath := *(flag.String("config",
				"/home/anton/url-shortener/config/local.yaml",
				"Path to the configuration file"))
	flag.Parse()

	if _, err := os.Stat(configPath); os.IsExist(err) {
		slog.Error("config file does not exist: ", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("cannot read config", err)
	}

	return  &cfg
}