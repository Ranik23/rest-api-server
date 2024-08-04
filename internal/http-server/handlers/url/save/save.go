package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	//"url-shortener/internal/lib/random"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)



type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string  `json:"alias,omitempty"`
}

//const aliasLength = 4

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// Handler Function for saving a user.
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op),

			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", err)
			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body is decoded", slog.Any("request", req))

		if err := validator.New().Struct(&req); err != nil {
			validateError := err.(validator.ValidationErrors)

			log.Error("invalid request", err)
			render.JSON(w, r, resp.ValidationError(validateError))

			return
		}

		//TODO

		// alias := req.Alias

		// if alias == "" {
		// 	alias = random
		// }

		id, err := urlSaver.SaveURL(req.URL, req.Alias); _ = id

		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
	}
}
