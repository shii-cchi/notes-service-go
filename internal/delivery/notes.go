package delivery

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"notes-service-go/internal/delivery/dto"
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
		r.Post("/", h.createHandler)
	})

	return rg
}

func (h NotesHandler) getHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	notes, err := h.notesService.GetNotes(accessToken)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrInvalidAccessToken) {
			respondWithError(w, http.StatusUnauthorized, domain.ErrInvalidAccessToken)
			return
		}
		respondWithError(w, http.StatusInternalServerError, domain.ErrGettingNotes)
		return
	}

	respondWithJSON(w, http.StatusOK, notes)
}

func (h NotesHandler) createHandler(w http.ResponseWriter, r *http.Request) {
	noteInput := dto.NoteInputDto{}
	if err := json.NewDecoder(r.Body).Decode(&noteInput); err != nil {
		log.Printf(domain.ErrParsingNoteInput+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, domain.ErrParsingNoteInput)
		return
	}

	if err := h.validator.Struct(&noteInput); err != nil {
		log.Printf(domain.ErrInvalidNoteInput+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, domain.ErrInvalidNoteInput)
		return
	}

	accessToken := r.Header.Get("Authorization")

	note, err := h.notesService.CreateNote(noteInput, accessToken)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrInvalidAccessToken) {
			respondWithError(w, http.StatusUnauthorized, domain.ErrInvalidAccessToken)
			return
		}
		if strings.HasPrefix(err.Error(), domain.ErrSpellingText) {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, domain.ErrCreatingNote)
		return
	}

	respondWithJSON(w, http.StatusCreated, note)
}
