package get

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/storage"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)



type URLGetter interface {
	GetURL(alias string) (string, error)
}

type Request struct {
	alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	url string `json:"url"`
}

// Handler Function For Getting A User
func New(log *slog.Logger, urlgetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers.url.get"

		log := log.With(
			slog.String("op", op),

			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("could not decode", sl.Err(err)) // TODO: err needs to be converted to slog.Attr
			return
		}

		url, err := urlgetter.GetURL(req.alias)

		if err != nil {

			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("alias not found", slog.String("alias",req.alias))
				render.JSON(w, r, resp.Error("alias not found"))
			}

			log.Error("could not find the url", sl.Err(err)) // TODO: err needs to be converted to slog.Attr

			return
		}


		render.JSON(w, r, Response {
			Response: resp.OK(),
			url : url,
		})
	}
}
