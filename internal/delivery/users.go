package delivery

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"notes-service-go/internal/delivery/dto"
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
		log.Printf(domain.ErrParsingUserCredentials+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, domain.ErrParsingUserCredentials)
		return
	}

	if err := h.validator.Struct(&userCredentials); err != nil {
		log.Printf(domain.ErrInvalidUserCredentials+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, domain.ErrInvalidUserCredentials)
		return
	}

	user, refreshToken, err := h.usersService.CreateUser(userCredentials)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrUserAlreadyExists) {
			respondWithError(w, http.StatusBadRequest, domain.ErrUserAlreadyExists)
			return
		}
		respondWithError(w, http.StatusInternalServerError, domain.ErrCreatingUser)
		return
	}
	setCookie(w, refreshToken, h.refreshTokenTTL)
	respondWithJSON(w, http.StatusCreated, user)
}

func (h UsersHandler) refreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			respondWithError(w, http.StatusUnauthorized, domain.ErrCookieNotFound)
			return
		}
		respondWithError(w, http.StatusUnauthorized, domain.ErrGettingRefreshTokenFromCookie)
		return
	}
	refreshToken := cookie.Value

	accessToken := r.Header.Get("Authorization")

	user, refreshToken, err := h.usersService.Refresh(refreshToken, accessToken)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrInvalidAccessToken) {
			respondWithError(w, http.StatusUnauthorized, domain.ErrInvalidAccessToken)
			return
		}
		if strings.HasPrefix(err.Error(), domain.ErrInvalidRefreshToken) {
			respondWithError(w, http.StatusUnauthorized, domain.ErrInvalidRefreshToken)
			return
		}
		respondWithError(w, http.StatusInternalServerError, domain.ErrRefresh)
		return
	}
	setCookie(w, refreshToken, h.refreshTokenTTL)
	respondWithJSON(w, http.StatusOK, user)
}

func (h UsersHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	userCredentials := dto.UserCredentialsDto{}
	if err := json.NewDecoder(r.Body).Decode(&userCredentials); err != nil {
		log.Printf(domain.ErrParsingUserCredentials+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, domain.ErrParsingUserCredentials)
		return
	}

	if err := h.validator.Struct(&userCredentials); err != nil {
		log.Printf(domain.ErrInvalidUserCredentials+" :%s\n", err)
		respondWithError(w, http.StatusBadRequest, domain.ErrInvalidUserCredentials)
		return
	}

	user, refreshToken, err := h.usersService.Login(userCredentials)
	if err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrUserNotFound) || strings.HasPrefix(err.Error(), domain.ErrWrongCredentials) {
			respondWithError(w, http.StatusBadRequest, domain.ErrWrongCredentials)
			return
		}
		respondWithError(w, http.StatusInternalServerError, domain.ErrLogin)
		return
	}
	setCookie(w, refreshToken, h.refreshTokenTTL)
	respondWithJSON(w, http.StatusOK, user)
}

func (h UsersHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	if err := h.usersService.Logout(accessToken); err != nil {
		log.Println(err)
		if strings.HasPrefix(err.Error(), domain.ErrInvalidAccessToken) {
			respondWithError(w, http.StatusUnauthorized, domain.ErrInvalidAccessToken)
			return
		}
		respondWithError(w, http.StatusInternalServerError, domain.ErrLogout)
		return
	}

	deleteCookie(w)
}
