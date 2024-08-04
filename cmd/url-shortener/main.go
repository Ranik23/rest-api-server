package main

import (
	//"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/get"
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

	router.Handle("/save", save.New(log, dataBase))
	router.Handle("/get", get.New(log, dataBase))


	srv := &http.Server{
		Addr: cfg.Address,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to establish connection to the server", sl.Err(err))
			return
		}
	}()


	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	s := <-c; _ = s 

	log.Info("app interrupted by a signal")
}