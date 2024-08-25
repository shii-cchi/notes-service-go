package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"notes-service-go/pkg/hash"
	"time"
)

const (
	ErrUnexpectedSigningMethod = "unexpected signing method"
	ErrGettingClaims           = "error getting user claims from token"
)

type TokenManager interface {
	NewAccessToken(userID string, accessTokenTTL time.Duration) (string, error)
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

func (m *Manager) NewAccessToken(userID string, accessTokenTTL time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   userID,
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

	return refreshTokenStr, hashedRefreshToken, nil
}

func (m *Manager) ParseAccessToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(ErrUnexpectedSigningMethod+": %v", token.Header["alg"])
		}

		return []byte(m.accessSigningKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	}

	return "", fmt.Errorf(ErrGettingClaims)
}

func (m *Manager) IsValidRefreshToken(hashedRefreshToken, refreshToken string) bool {
	return m.hasher.IsValidData(hashedRefreshToken, refreshToken)
}
