package dto

import (
	_ "github.com/go-playground/validator/v10"
)

type UserCredentialsDto struct {
	Login    string `json:"login" validate:"required,min=6"`
	Password string `json:"password" validate:"required,min=6"`
}
