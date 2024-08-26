package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"notes-service-go/internal/constants"
	"notes-service-go/internal/database"
	"notes-service-go/internal/delivery/dto"
	"notes-service-go/pkg/auth"
	"notes-service-go/pkg/spell"
)

type NotesService struct {
	Repo         *database.Queries
	Speller      spell.Speller
	TokenManager auth.TokenManager
}

func NewNotesService(repo *database.Queries, speller spell.Speller, tokenManager auth.TokenManager) *NotesService {
	return &NotesService{
		Repo:         repo,
		Speller:      speller,
		TokenManager: tokenManager,
	}
}

func (s *NotesService) GetNotes(accessToken string) ([]dto.NoteResponseDto, error) {
	userIDStr, err := s.TokenManager.ParseAccessToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrInvalidAccessToken+" :%s\n", err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrParsingID+" :%s\n", err)
	}

	notes, err := s.Repo.GetNotes(context.Background(), userID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrGettingNotes+" :%s\n", err)
	}

	return s.newNotesResponseDto(notes), nil
}

func (s *NotesService) CreateNote(noteInput dto.NoteInputDto, accessToken string) (dto.NoteResponseDto, error) {
	userIDStr, err := s.TokenManager.ParseAccessToken(accessToken)
	if err != nil {
		return dto.NoteResponseDto{}, fmt.Errorf(constants.ErrInvalidAccessToken+" :%s\n", err)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return dto.NoteResponseDto{}, fmt.Errorf(constants.ErrParsingID+" :%s\n", err)
	}

	spellingErrors, err := s.Speller.CheckText(noteInput.Content)
	if err != nil {
		return dto.NoteResponseDto{}, fmt.Errorf(constants.ErrCheckingSpellingErrors+" :%s\n", err)
	}

	if len(spellingErrors) != 0 {
		return dto.NoteResponseDto{}, errors.New(constants.ErrSpellingText + ". " + s.Speller.FormatErrors(spellingErrors))
	}

	note, err := s.Repo.CreateNote(context.Background(), database.CreateNoteParams{Name: noteInput.Name, Content: noteInput.Content, UserID: userID})
	if err != nil {
		return dto.NoteResponseDto{}, fmt.Errorf(constants.ErrCreatingNote+" :%s\n", err)
	}

	return s.newNoteResponseDto(note), nil
}

func (s *NotesService) newNoteResponseDto(note database.CreateNoteRow) dto.NoteResponseDto {
	return dto.NoteResponseDto{
		ID:      note.ID,
		Name:    note.Name,
		Content: note.Content,
	}
}

func (s *NotesService) newNotesResponseDto(notes []database.GetNotesRow) []dto.NoteResponseDto {
	dtos := make([]dto.NoteResponseDto, len(notes))
	for i, note := range notes {
		dtos[i] = dto.NoteResponseDto{
			ID:      note.ID,
			Name:    note.Name,
			Content: note.Content,
		}
	}
	return dtos
}
