package main

import (
	"fmt"
	"log/slog"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/logger"
	"url-shortener/internal/storage/sqlite"
)






func main() {

	cfg := config.MustLoad()
	fmt.Println(cfg)
	log := logger.SetUpLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	dataBase, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init database", sl.Err(err))
	}

	
}