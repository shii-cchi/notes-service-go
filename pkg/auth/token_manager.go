package auth

import (
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
	NewAccessToken(userID uuid.UUID) (string, error)
	NewRefreshToken(userID uuid.UUID) (string, error)
	ParseAccessToken(accessToken string) (string, error)
	ParseRefreshToken(refreshToken string) (string, error)
}

type Manager struct {
	accessTTL         time.Duration
	refreshTTL        time.Duration
	accessSigningKey  string
	refreshSigningKey string
	hasher            hash.Hasher
}

func NewManager(accessTTL time.Duration, refreshTTL time.Duration, accessSigningKey string, refreshSigningKey string, hasher hash.Hasher) *Manager {
	return &Manager{
		accessTTL:         accessTTL,
		refreshTTL:        refreshTTL,
		accessSigningKey:  accessSigningKey,
		refreshSigningKey: refreshSigningKey,
		hasher:            hasher,
	}
}

func (m *Manager) newToken(userID uuid.UUID, ttl time.Duration, signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   userID.String(),
	})

	return token.SignedString([]byte(signingKey))
}

func (m *Manager) NewAccessToken(userID uuid.UUID) (string, error) {
	return m.newToken(userID, m.accessTTL, m.accessSigningKey)
}

func (m *Manager) NewRefreshToken(userID uuid.UUID) (string, error) {
	return m.newToken(userID, m.refreshTTL, m.refreshSigningKey)
}

func (m *Manager) parseToken(receivedToken string, signingKey string) (string, error) {
	token, err := jwt.Parse(receivedToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(errUnexpectedSigningMethod+": %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["sub"].(string), nil
	}

	return "", fmt.Errorf(errGettingClaims)
}

func (m *Manager) ParseAccessToken(accessToken string) (string, error) {
	if accessToken == "" {
		return "", errors.New(errAccessTokenUndefined)
	}

	if strings.HasPrefix(accessToken, accessTokenPrefix) {
		accessToken = strings.TrimPrefix(accessToken, accessTokenPrefix)
	}

	return m.parseToken(accessToken, m.accessSigningKey)
}

func (m *Manager) ParseRefreshToken(refreshToken string) (string, error) {
	return m.parseToken(refreshToken, m.refreshSigningKey)
}
