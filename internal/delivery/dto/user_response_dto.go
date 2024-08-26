package dto

import (
	"github.com/google/uuid"
)

type UserResponseDto struct {
	ID          uuid.UUID `json:"id"`
	Login       string    `json:"login"`
	AccessToken string    `json:"access_token"`
}
