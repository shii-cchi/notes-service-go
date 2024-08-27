package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"notes-service-go/internal/delivery"
	"notes-service-go/internal/delivery/dto"
	"notes-service-go/internal/delivery/middleware"
	"notes-service-go/internal/domain"
	"notes-service-go/internal/service"
	"strings"
	"time"
)

type UsersHandler struct {
	usersService service.Users
	validator    *validator.Validate

	refreshTokenTTL time.Duration
}

func NewUsersHandler(usersService service.Users, validator *validator.Validate, refreshTokenTTL time.Duration) *UsersHandler {
	return &UsersHandler{
		usersService:    usersService,
		validator:       validator,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (h UsersHandler) usersHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Post("/register", middleware.CheckUserCredentialsInput(h.validator, h.registerHandler))
		r.Get("/refresh", h.refreshHandler)
		r.Post("/login", middleware.CheckUserCredentialsInput(h.validator, h.loginHandler))
		r.Get("/logout", h.logoutHandler)
	})

	return rg
}

func (h UsersHandler) registerHandler(w http.ResponseWriter, r *http.Request, userCredentials dto.UserCredentialsDto) {
	user, refreshToken, err := h.usersService.CreateUser(userCredentials)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrUserAlreadyExists) {
			delivery.RespondWithError(w, http.StatusBadRequest, domain.ErrUserAlreadyExists)
			return
		}
		delivery.RespondWithError(w, http.StatusInternalServerError, domain.ErrCreatingUser)
		return
	}
	delivery.SetCookie(w, refreshToken, h.refreshTokenTTL)
	delivery.RespondWithJSON(w, http.StatusCreated, user)
}

func (h UsersHandler) refreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Println(err)
		if err == http.ErrNoCookie {
			delivery.RespondWithError(w, http.StatusUnauthorized, domain.ErrRefreshTokenUndefined)
			return
		}
		delivery.RespondWithError(w, http.StatusUnauthorized, domain.ErrGettingRefreshTokenFromCookie)
		return
	}
	refreshToken := cookie.Value

	user, refreshToken, err := h.usersService.Refresh(refreshToken)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrInvalidAccessToken) || strings.HasPrefix(err.Error(), domain.ErrAccessTokenUndefined) {
			delivery.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if strings.HasPrefix(err.Error(), domain.ErrInvalidRefreshToken) {
			delivery.RespondWithError(w, http.StatusUnauthorized, domain.ErrInvalidRefreshToken)
			return
		}
		delivery.RespondWithError(w, http.StatusInternalServerError, domain.ErrRefresh)
		return
	}
	delivery.SetCookie(w, refreshToken, h.refreshTokenTTL)
	delivery.RespondWithJSON(w, http.StatusOK, user)
}

func (h UsersHandler) loginHandler(w http.ResponseWriter, r *http.Request, userCredentials dto.UserCredentialsDto) {
	user, refreshToken, err := h.usersService.Login(userCredentials)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrUserNotFound) || strings.HasPrefix(err.Error(), domain.ErrWrongPassword) {
			delivery.RespondWithError(w, http.StatusBadRequest, domain.ErrWrongCredentials)
			return
		}
		delivery.RespondWithError(w, http.StatusInternalServerError, domain.ErrLogin)
		return
	}
	delivery.SetCookie(w, refreshToken, h.refreshTokenTTL)
	delivery.RespondWithJSON(w, http.StatusOK, user)
}

func (h UsersHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	if err := h.usersService.Logout(accessToken); err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrInvalidAccessToken) {
			delivery.RespondWithError(w, http.StatusUnauthorized, domain.ErrInvalidAccessToken)
			return
		}
		delivery.RespondWithError(w, http.StatusInternalServerError, domain.ErrLogout)
		return
	}

	delivery.DeleteCookie(w)
}
