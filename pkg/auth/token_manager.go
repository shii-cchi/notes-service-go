package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"notes-service-go/pkg/hash"
	"strings"
	"time"
)

const (
	errUnexpectedSigningMethod = "unexpected signing method"
	errAccessTokenUndefined    = "access token is undefined"
	errGettingClaims           = "error getting user claims from token"

	accessTokenPrefix = "Bearer "
)

type TokenManager interface {
	NewAccessToken(userID uuid.UUID, accessTokenTTL time.Duration) (string, error)
	NewRefreshToken() (string, string, error)
	ParseAccessToken(accessToken string) (string, error)
	IsValidRefreshToken(hashedRefreshToken, refreshToken string) bool
}

type Manager struct {
	accessSigningKey string
	hasher           hash.Hasher
}

func NewManager(accessSigningKey string, hasher hash.Hasher) *Manager {
	return &Manager{
		accessSigningKey: accessSigningKey,
		hasher:           hasher,
	}
}

func (m *Manager) NewAccessToken(userID uuid.UUID, accessTokenTTL time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   userID.String(),
	})

	return token.SignedString([]byte(m.accessSigningKey))
}

func (m *Manager) NewRefreshToken() (string, string, error) {
	refreshToken := make([]byte, 32)

	_, err := rand.Read(refreshToken)
	if err != nil {
		return "", "", err
	}

	refreshTokenStr := base64.StdEncoding.EncodeToString(refreshToken)

	hashedRefreshToken, err := m.hasher.Hash(refreshTokenStr)
	if err != nil {
		return "", "", err
	}

	return refreshTokenStr, hashedRefreshToken, nil
}

func (m *Manager) ParseAccessToken(accessToken string) (string, error) {
	if accessToken == "" {
		return "", errors.New(errAccessTokenUndefined)
	}

	if strings.HasPrefix(accessToken, accessTokenPrefix) {
		accessToken = strings.TrimPrefix(accessToken, accessTokenPrefix)
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(errUnexpectedSigningMethod+": %v", token.Header["alg"])
		}

		return []byte(m.accessSigningKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	}

	return "", fmt.Errorf(errGettingClaims)
}

func (m *Manager) IsValidRefreshToken(hashedRefreshToken, refreshToken string) bool {
	return m.hasher.IsValidData(hashedRefreshToken, refreshToken)
}
