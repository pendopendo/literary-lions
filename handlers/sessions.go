package handlers

import (
	"database/sql"
	"literary-lions/models"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateSession(w http.ResponseWriter, r *http.Request, db *sql.DB, userID int) error {
	sessionToken := uuid.New().String()
	expiration := time.Now().Add(24 * time.Hour)

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiration,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	return models.SaveSession(db, sessionToken, userID, expiration)
}

func getUserIDFromSession(r *http.Request) int {
	userID := r.Context().Value("userID")
	if userID != nil {
		return userID.(int)
	}
	return 0
}

func ClearSession(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}
