package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	//"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	//"github.com/unrolled/render"
	"github.com/go-chi/render"
)



type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string  `json:"alias,omitempty"`
}

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

			render.JSON(w, r, resp.Error("invalid request"))
			render.JSON(w, r, resp.ValidationError(validateError))

			return
		}

		code, err := urlSaver.SaveURL(req.URL, req.URL); _ = code

		if err != nil {
			log.Error("failed to save the user", err)
			return
		}
	}
}
