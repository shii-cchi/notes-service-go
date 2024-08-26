package delivery

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"notes-service-go/internal/constants"
	"notes-service-go/internal/delivery/dto"
	"notes-service-go/internal/service"
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
		r.Get("/register", h.registerHandler)
		r.Get("/refresh", h.refreshHandler)
		r.Get("/login", h.loginHandler)
		r.Get("/logout", h.logoutHandler)
	})

	return rg
}

func (h UsersHandler) registerHandler(w http.ResponseWriter, r *http.Request) {
	userCredentials := dto.UserCredentialsDto{}
	if err := json.NewDecoder(r.Body).Decode(&userCredentials); err != nil {
		log.Printf(constants.ErrParsingUserCredentials+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, constants.ErrParsingUserCredentials)
		return
	}

	if err := h.validator.Struct(&userCredentials); err != nil {
		log.Printf(constants.ErrInvalidUserCredentials+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, constants.ErrInvalidUserCredentials+" "+constants.UserCredentialsUsage)
		return
	}

	user, refreshToken, err := h.usersService.CreateUser(userCredentials)
	if err != nil {
		// TODO handle error 400 - if user already exists, 500 - another
	}
	setCookie(w, refreshToken, h.refreshTokenTTL)
	respondWithJSON(w, http.StatusCreated, user)
}

func (h UsersHandler) refreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			respondWithError(w, http.StatusUnauthorized, constants.ErrCookieNotFound)
			return
		}
		respondWithError(w, http.StatusUnauthorized, constants.ErrGettingRefreshToken)
		return
	}
	refreshToken := cookie.Value

	accessToken := r.Header.Get("Authorization")

	user, refreshToken, err := h.usersService.Refresh(refreshToken, accessToken)
	if err != nil {
		// TODO handle error 401 - unauth, 500 - another
	}
	setCookie(w, refreshToken, h.refreshTokenTTL)
	respondWithJSON(w, http.StatusOK, user)
}

func (h UsersHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	userCredentials := dto.UserCredentialsDto{}
	if err := json.NewDecoder(r.Body).Decode(&userCredentials); err != nil {
		log.Printf(constants.ErrParsingUserCredentials+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, constants.ErrParsingUserCredentials)
		return
	}

	if err := h.validator.Struct(&userCredentials); err != nil {
		log.Printf(constants.ErrInvalidUserCredentials+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, constants.ErrInvalidUserCredentials+" "+constants.UserCredentialsUsage)
		return
	}

	user, refreshToken, err := h.usersService.Login(userCredentials)
	if err != nil {
		// TODO handle error 500 - another
	}
	setCookie(w, refreshToken, h.refreshTokenTTL)
	respondWithJSON(w, http.StatusOK, user)
}

func (h UsersHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	if err := h.usersService.Logout(accessToken); err != nil {
		// TODO handle error 500 - another, 401 - unauth
	}

	deleteCookie(w)
}
