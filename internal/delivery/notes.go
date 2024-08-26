package delivery

import (
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"net/http"
	"notes-service-go/internal/service"
)

type NotesHandler struct {
	notesService service.Notes
	validator    *validator.Validate
}

func NewNoteHandler(notesService service.Notes, validator *validator.Validate) *NotesHandler {
	return &NotesHandler{
		notesService: notesService,
		validator:    validator,
	}
}

func (h NotesHandler) notesHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", h.getHandler)
		r.Post("/", h.createHandler)
	})

	return rg
}

func (h NotesHandler) getHandler(w http.ResponseWriter, r *http.Request) {

}

func (h NotesHandler) createHandler(w http.ResponseWriter, r *http.Request) {

}
