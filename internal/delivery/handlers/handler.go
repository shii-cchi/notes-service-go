package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"notes-service-go/internal/service"
	"time"
)

type Handler struct {
	UsersHandler *UsersHandler
	NotesHandler *NotesHandler
}

func NewHandler(services *service.Services, validator *validator.Validate, refreshTokenTTL time.Duration) *Handler {
	return &Handler{
		UsersHandler: NewUsersHandler(services.Users, validator, refreshTokenTTL),
		NotesHandler: NewNoteHandler(services.Notes, validator),
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Mount("/users", h.UsersHandler.usersHandlers())
	r.Mount("/notes", h.NotesHandler.notesHandlers())
}
