package middleware

import (
	"context"
	"database/sql"
	"literary-lions-forum/models"
	"net/http"
)

func AuthMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, err := models.GetSessionByToken(db, cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserIDFromSession(r *http.Request) int {
	userID := r.Context().Value("userID")
	if userID != nil {
		return userID.(int)
	}
	return 0
}
