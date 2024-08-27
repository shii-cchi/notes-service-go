package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"notes-service-go/internal/database"
	"notes-service-go/internal/delivery/dto"
	"notes-service-go/internal/domain"
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

func NewUsersService(repo *database.Queries, hasher hash.Hasher, tokenManager auth.TokenManager) *UsersService {
	return &UsersService{
		Repo:         repo,
		Hasher:       hasher,
		TokenManager: tokenManager,
	}
}

func (s *UsersService) CreateUser(userCredentials dto.UserCredentialsDto) (dto.UserResponseDto, string, error) {
	exist, err := s.Repo.CheckUserExist(context.Background(), userCredentials.Login)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrCheckingUserExist+": %s\n", err)
	}
	if exist {
		return dto.UserResponseDto{}, "", errors.New(domain.ErrUserAlreadyExists)
	}

	hashedPassword, err := s.Hasher.Hash(userCredentials.Password)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrHashingPassword+": %s\n", err)
	}

	userID, err := s.Repo.CreateUser(context.Background(), database.CreateUserParams{Login: userCredentials.Login, Password: hashedPassword})
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrCreatingUser+": %s\n", err)
	}

	refreshToken, err := s.TokenManager.NewRefreshToken(userID)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrCreatingRefreshToken+": %s\n", err)
	}

	accessToken, err := s.TokenManager.NewAccessToken(userID)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrCreatingAccessToken+": %s\n", err)
	}

	if err = s.Repo.SaveRefreshToken(context.Background(), database.SaveRefreshTokenParams{ID: userID, RefreshToken: refreshToken}); err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrSavingRefreshToken+": %s\n", err)
	}

	return dto.UserResponseDto{ID: userID, AccessToken: accessToken}, refreshToken, nil
}

func (s *UsersService) Refresh(refreshToken string) (dto.UserResponseDto, string, error) {
	userIDStr, err := s.TokenManager.ParseAccessToken(refreshToken)
	if err != nil {
		if err.Error() == domain.ErrRefreshTokenUndefined {
			return dto.UserResponseDto{}, "", err
		}
		return dto.UserResponseDto{}, "", errors.New(domain.ErrInvalidRefreshToken)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrParsingID+" :%s\n", err)
	}

	storedRefreshToken, err := s.Repo.GetRefreshTokenById(context.Background(), userID)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrGettingRefreshTokenFromDB+" :%s\n", err)
	}

	if storedRefreshToken != refreshToken {
		return dto.UserResponseDto{}, "", errors.New(domain.ErrInvalidRefreshToken)
	}

	refreshToken, err = s.TokenManager.NewRefreshToken(userID)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrCreatingRefreshToken+": %s\n", err)
	}

	accessToken, err := s.TokenManager.NewAccessToken(userID)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrCreatingAccessToken+": %s\n", err)
	}

	if err = s.Repo.SaveRefreshToken(context.Background(), database.SaveRefreshTokenParams{ID: userID, RefreshToken: refreshToken}); err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrSavingRefreshToken+": %s\n", err)
	}

	return dto.UserResponseDto{ID: userID, AccessToken: accessToken}, refreshToken, nil
}

func (s *UsersService) Login(userCredentials dto.UserCredentialsDto) (dto.UserResponseDto, string, error) {
	user, err := s.Repo.GetUserByLogin(context.Background(), userCredentials.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrUserNotFound+": %s\n", err)
		}
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrGettingPassword+": %s\n", err)
	}

	valid := s.Hasher.IsValidData(user.Password, userCredentials.Password)
	if !valid {
		return dto.UserResponseDto{}, "", errors.New(domain.ErrWrongPassword)
	}

	refreshToken, err := s.TokenManager.NewRefreshToken(user.ID)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrCreatingRefreshToken+": %s\n", err)
	}

	accessToken, err := s.TokenManager.NewAccessToken(user.ID)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrCreatingAccessToken+": %s\n", err)
	}

	if err = s.Repo.SaveRefreshToken(context.Background(), database.SaveRefreshTokenParams{ID: user.ID, RefreshToken: refreshToken}); err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(domain.ErrSavingRefreshToken+": %s\n", err)
	}

	return dto.UserResponseDto{ID: user.ID, AccessToken: accessToken}, refreshToken, nil
}

func (s *UsersService) Logout(accessToken string) error {
	userIDStr, err := s.TokenManager.ParseAccessToken(accessToken)
	if err != nil {
		if err.Error() == domain.ErrAccessTokenUndefined {
			return err
		}
		return errors.New(domain.ErrInvalidAccessToken)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf(domain.ErrParsingID+" :%s\n", err)
	}

	if err = s.Repo.Logout(context.Background(), userID); err != nil {
		return fmt.Errorf(domain.ErrLogout+" :%s\n", err)
	}

	return nil
}
