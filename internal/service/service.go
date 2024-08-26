package service

import (
	"notes-service-go/internal/database"
	"notes-service-go/internal/delivery/dto"
	"notes-service-go/pkg/auth"
	"notes-service-go/pkg/hash"
	"time"
)

type Users interface {
	CreateUser(userCredentials dto.UserCredentialsDto) (dto.UserResponseDto, string, error)
	Refresh(refreshToken string, accessToken string) (dto.UserResponseDto, string, error)
	Login(userCredentials dto.UserCredentialsDto) (dto.UserResponseDto, string, error)
	Logout(accessToken string) error
}

type Notes interface {
	GetNotes(accessToken string) ([]dto.NoteResponseDto, error)
	CreateNote(noteInput dto.NoteInputDto, accessToken string) (dto.NoteResponseDto, error)
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
	notesService := NewNotesService(deps.Repo, deps.TokenManager)

	return &Services{
		Users: usersService,
		Notes: notesService,
	}
}
