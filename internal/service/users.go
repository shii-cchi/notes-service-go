package service

import (
	"notes-service-go/internal/database"
	"notes-service-go/internal/delivery/dto"
	"notes-service-go/pkg/auth"
	"notes-service-go/pkg/hash"
	"time"
)

type UsersService struct {
	Repo         *database.Queries
	Hasher       hash.Hasher
	TokenManager auth.TokenManager

	AccessTokenTTL time.Duration
}

func NewUsersService(repo *database.Queries, hasher hash.Hasher, tokenManager auth.TokenManager, accessTokenTTL time.Duration) *UsersService {
	return &UsersService{
		Repo:           repo,
		Hasher:         hasher,
		TokenManager:   tokenManager,
		AccessTokenTTL: accessTokenTTL,
	}
}

func (s *UsersService) CreateUser(userCredentials dto.UserCredentialsDto) (dto.UserResponseDto, string, error) {

}

func (s *UsersService) Refresh(refreshToken string, accessToken string) (dto.UserResponseDto, string, error) {

}

func (s *UsersService) Login(userCredentials dto.UserCredentialsDto) (dto.UserResponseDto, string, error) {

}

func (s *UsersService) Logout(accessToken string) error {

}
