package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"notes-service-go/internal/delivery"
	"notes-service-go/internal/delivery/dto"
	"notes-service-go/internal/delivery/middleware"
	"notes-service-go/internal/domain"
	"notes-service-go/internal/service"
	"strings"
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
		r.Post("/", middleware.CheckNoteInput(h.validator, h.createHandler))
	})

	return rg
}

func (h NotesHandler) getHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	notes, err := h.notesService.GetNotes(accessToken)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrInvalidAccessToken) || strings.HasPrefix(err.Error(), domain.ErrAccessTokenUndefined) {
			delivery.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		delivery.RespondWithError(w, http.StatusInternalServerError, domain.ErrGettingNotes)
		return
	}

	delivery.RespondWithJSON(w, http.StatusOK, notes)
}

func (h NotesHandler) createHandler(w http.ResponseWriter, r *http.Request, noteInput dto.NoteInputDto) {
	accessToken := r.Header.Get("Authorization")

	note, err := h.notesService.CreateNote(noteInput, accessToken)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrInvalidAccessToken) || strings.HasPrefix(err.Error(), domain.ErrAccessTokenUndefined) {
			delivery.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if strings.HasPrefix(err.Error(), domain.ErrSpellingText) {
			delivery.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		delivery.RespondWithError(w, http.StatusInternalServerError, domain.ErrCreatingNote)
		return
	}

	delivery.RespondWithJSON(w, http.StatusCreated, note)
}
