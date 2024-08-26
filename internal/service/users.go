package service

import (
	"context"
	"errors"
	"fmt"
	"notes-service-go/internal/constants"
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
	exist, err := s.Repo.CheckUserExist(context.Background(), userCredentials.Login)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrCheckingUserExist+": %s\n", err)
	}
	if exist {
		return dto.UserResponseDto{}, "", errors.New(constants.ErrUserAlreadyExists)
	}

	hashedPassword, err := s.Hasher.Hash(userCredentials.Password)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrHashingPassword+": %s\n", err)
	}

	userID, err := s.Repo.CreateUser(context.Background(), database.CreateUserParams{Login: userCredentials.Login, Password: hashedPassword})
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrCreatingUser+": %s\n", err)
	}

	refreshToken, hashedRefreshToken, err := s.TokenManager.NewRefreshToken()
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrCreatingRefreshToken+": %s\n", err)
	}

	accessToken, err := s.TokenManager.NewAccessToken(userID, s.AccessTokenTTL)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrCreatingAccessToken+": %s\n", err)
	}

	if err = s.Repo.SaveRefreshToken(context.Background(), database.SaveRefreshTokenParams{ID: userID, RefreshToken: hashedRefreshToken}); err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrSavingRefreshToken+": %s\n", err)
	}

	return dto.UserResponseDto{ID: userID, AccessToken: accessToken}, refreshToken, nil
}

func (s *UsersService) Refresh(refreshToken string, accessToken string) (dto.UserResponseDto, string, error) {
	userID, err := s.TokenManager.ParseAccessToken(accessToken)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrInvalidAccessToken+" :%s\n", err)
	}

	hashedStoredRefreshToken, err := s.Repo.GetRefreshTokenById(context.Background(), userID)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrGettingRefreshTokenFromDB+" :%s\n", err)
	}

	valid := s.Hasher.IsValidData(hashedStoredRefreshToken, refreshToken)
	if !valid {
		return dto.UserResponseDto{}, "", errors.New(constants.ErrInvalidRefreshToken)
	}

	refreshToken, hashedRefreshToken, err := s.TokenManager.NewRefreshToken()
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrCreatingRefreshToken+": %s\n", err)
	}

	accessToken, err = s.TokenManager.NewAccessToken(userID, s.AccessTokenTTL)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrCreatingAccessToken+": %s\n", err)
	}

	if err = s.Repo.SaveRefreshToken(context.Background(), database.SaveRefreshTokenParams{ID: userID, RefreshToken: hashedRefreshToken}); err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrSavingRefreshToken+": %s\n", err)
	}

	return dto.UserResponseDto{ID: userID, AccessToken: accessToken}, refreshToken, nil
}

func (s *UsersService) Login(userCredentials dto.UserCredentialsDto) (dto.UserResponseDto, string, error) {
	user, err := s.Repo.GetUserByLogin(context.Background(), userCredentials.Login)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrGettingPassword+": %s\n", err)
	}
	if user.Password == "" {
		return dto.UserResponseDto{}, "", errors.New(constants.ErrUserNotFound)
	}

	valid := s.Hasher.IsValidData(user.Password, userCredentials.Password)
	if !valid {
		return dto.UserResponseDto{}, "", errors.New(constants.ErrWrongCredentials)
	}

	refreshToken, hashedRefreshToken, err := s.TokenManager.NewRefreshToken()
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrCreatingRefreshToken+": %s\n", err)
	}

	accessToken, err := s.TokenManager.NewAccessToken(user.ID, s.AccessTokenTTL)
	if err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrCreatingAccessToken+": %s\n", err)
	}

	if err = s.Repo.SaveRefreshToken(context.Background(), database.SaveRefreshTokenParams{ID: user.ID, RefreshToken: hashedRefreshToken}); err != nil {
		return dto.UserResponseDto{}, "", fmt.Errorf(constants.ErrSavingRefreshToken+": %s\n", err)
	}

	return dto.UserResponseDto{ID: user.ID, AccessToken: accessToken}, refreshToken, nil
}

func (s *UsersService) Logout(accessToken string) error {
	userID, err := s.TokenManager.ParseAccessToken(accessToken)
	if err != nil {
		return fmt.Errorf(constants.ErrInvalidAccessToken+" :%s\n", err)
	}

	if err = s.Repo.Logout(context.Background(), userID); err != nil {
		return fmt.Errorf(constants.ErrLogout+" :%s\n", err)
	}

	return nil
}
