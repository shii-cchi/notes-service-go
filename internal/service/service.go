package service

import (
	"notes-service-go/internal/database"
	"notes-service-go/pkg/auth"
	"notes-service-go/pkg/hash"
	"time"
)

type Users interface {
}

type Notes interface {
}

type Services struct {
	Users Users
	Notes Notes
}

type Deps struct {
	Repo         *database.Queries
	Hasher       hash.Hasher
	TokenManager auth.TokenManager

	AccessTokenTTL time.Duration
}

func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repo, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL)
	notesService := NewNotesService(deps.Repo)

	return &Services{
		Users: usersService,
		Notes: notesService,
	}
}
