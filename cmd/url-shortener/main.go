package main

import (
	//"fmt"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)






func main() {

	cfg := config.MustLoad()
	//fmt.Println(cfg)

	log := sl.SetUpLogger(cfg.Env)
	
	log = log.With(slog.String("env", cfg.Env))

	dataBase, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init database", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.Logger)
	router.Use(middleware.Recoverer)

	router.Handle("/", save.New(log, dataBase))


	http.ListenAndServe("localhost:8080", router)

	//router.Use(middleware.URLFormat)


}