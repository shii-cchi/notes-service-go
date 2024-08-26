package dto

type NoteInputDto struct {
	Name    string `json:"name" validate:"required,min=1"`
	Content string `json:"content" validate:"required,min=1"`
}
