package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gen4ralz/movie-app/movie-service/internal/controller/movie"
)

// Handler defines a movie handler.
type Handler struct {
	ctrl *movie.Controller
}

// New creates a new movie HTTP handler.
func New(ctrl *movie.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
	}
}

// GetMovieDetails handles Get /movie requests.
func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")

	details, err := h.ctrl.Get(req.Context(), id)
	if err != nil && errors.Is(err, movie.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(details)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}