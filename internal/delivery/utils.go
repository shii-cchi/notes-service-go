package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal JSON responce: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func setCookie(w http.ResponseWriter, refreshToken string, refreshTTL time.Duration) {
	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(refreshTTL),
		Path:     "/users",
	}

	http.SetCookie(w, &cookie)
}

func deleteCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/users",
	}

	http.SetCookie(w, &cookie)
}
