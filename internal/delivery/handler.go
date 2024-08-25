package delivery

import (
	"github.com/go-chi/chi"
	"notes-service-go/internal/service"
	"time"
)

type Handler struct {
	services *service.Services

	refreshTokenTTL time.Duration
}

func NewHandler(services *service.Services, refreshTokenTTL time.Duration) *Handler {
	return &Handler{
		services:        services,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {}
