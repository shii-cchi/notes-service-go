package middleware

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"notes-service-go/internal/delivery"
	"notes-service-go/internal/delivery/dto"
	"notes-service-go/internal/domain"
)

func CheckUserCredentialsInput(v *validator.Validate, next func(http.ResponseWriter, *http.Request, dto.UserCredentialsDto)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userCredentials := dto.UserCredentialsDto{}
		if err := json.NewDecoder(r.Body).Decode(&userCredentials); err != nil {
			log.Printf(domain.ErrParsingUserCredentialsInput+" :%s\n", err)
			delivery.RespondWithError(w, http.StatusBadRequest, domain.ErrParsingUserCredentialsInput)
			return
		}

		if err := v.Struct(&userCredentials); err != nil {
			log.Printf(domain.ErrInvalidUserCredentialsInput+" :%s\n", err)
			delivery.RespondWithError(w, http.StatusBadRequest, domain.ErrInvalidUserCredentialsInput)
			return
		}

		next(w, r, userCredentials)
	}
}

func CheckNoteInput(v *validator.Validate, next func(http.ResponseWriter, *http.Request, dto.NoteInputDto)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteInput := dto.NoteInputDto{}
		if err := json.NewDecoder(r.Body).Decode(&noteInput); err != nil {
			log.Printf(domain.ErrParsingNoteInput+" :%s\n", err)
			delivery.RespondWithError(w, http.StatusBadRequest, domain.ErrParsingNoteInput)
			return
		}

		if err := v.Struct(&noteInput); err != nil {
			log.Printf(domain.ErrInvalidNoteInput+" :%s\n", err)
			delivery.RespondWithError(w, http.StatusBadRequest, domain.ErrInvalidNoteInput)
			return
		}

		next(w, r, noteInput)
	}
}
