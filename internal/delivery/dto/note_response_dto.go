package dto

import "github.com/google/uuid"

type NoteResponseDto struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
}
