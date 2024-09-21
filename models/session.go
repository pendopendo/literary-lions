package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           int
	UserID       int
	SessionToken string
	ExpiresAt    time.Time
}

func CreateSession(db *sql.DB, userID int) (string, error) {
	sessionToken := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err := db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", userID, sessionToken, expiresAt)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func GetSessionByToken(db *sql.DB, token string) (*Session, error) {
	session := &Session{}
	err := db.QueryRow("SELECT id, user_id, session_token, expires_at FROM sessions WHERE session_token = ?", token).Scan(&session.ID, &session.UserID, &session.SessionToken, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, sql.ErrNoRows
	}

	return session, nil
}

func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE session_token = ?", token)
	return err
}

func SaveSession(db *sql.DB, token string, userID int, expiresAt time.Time) error {
	_, err := db.Exec("INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)", token, userID, expiresAt)
	return err
}

func GetUserIDBySession(db *sql.DB, token string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
