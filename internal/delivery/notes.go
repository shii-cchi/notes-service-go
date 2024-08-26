package delivery

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"notes-service-go/internal/constants"
	"notes-service-go/internal/delivery/dto"
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
		if strings.HasPrefix(err.Error(), constants.ErrInvalidAccessToken) {
			respondWithError(w, http.StatusUnauthorized, constants.ErrInvalidAccessToken)
			return
		}
		respondWithError(w, http.StatusInternalServerError, constants.ErrGettingNotes)
		return
	}

	respondWithJSON(w, http.StatusOK, notes)
}

func (h NotesHandler) createHandler(w http.ResponseWriter, r *http.Request) {
	noteInput := dto.NoteInputDto{}
	if err := json.NewDecoder(r.Body).Decode(&noteInput); err != nil {
		log.Printf(constants.ErrParsingNoteInput+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, constants.ErrParsingNoteInput)
		return
	}

	if err := h.validator.Struct(&noteInput); err != nil {
		log.Printf(constants.ErrInvalidNoteInput+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, constants.ErrInvalidNoteInput+" "+constants.NoteInputUsage)
		return
	}

	accessToken := r.Header.Get("Authorization")

	note, err := h.notesService.CreateNote(noteInput, accessToken)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), constants.ErrInvalidAccessToken) {
			respondWithError(w, http.StatusUnauthorized, constants.ErrInvalidAccessToken)
			return
		}
		if strings.HasPrefix(err.Error(), constants.ErrSpellingText) {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, constants.ErrCreatingNote)
		return
	}

	respondWithJSON(w, http.StatusCreated, note)
}
