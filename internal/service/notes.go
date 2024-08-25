package service

import "notes-service-go/internal/database"

type NotesService struct {
	Repo *database.Queries
}

func NewNotesService(repo *database.Queries) *NotesService {
	return &NotesService{
		Repo: repo,
	}
}
