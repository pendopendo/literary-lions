package models

import (
	"database/sql"
)

func LikePost(db *sql.DB, userID, postID int, like bool) error {
	_, err := db.Exec("INSERT OR REPLACE INTO likes (user_id, post_id, like) VALUES (?, ?, ?)", userID, postID, like)
	return err
}

func LikeComment(db *sql.DB, userID, commentID int, like bool) error {
	_, err := db.Exec("INSERT OR REPLACE INTO likes (user_id, comment_id, like) VALUES (?, ?, ?)", userID, commentID, like)
	return err
}

func GetLikesForPost(db *sql.DB, postID int) (int, int, error) {
	var likes, dislikes int
	err := db.QueryRow("SELECT COUNT(CASE WHEN like = 1 THEN 1 END), COUNT(CASE WHEN like = 0 THEN 1 END) FROM likes WHERE post_id = ?", postID).Scan(&likes, &dislikes)
	return likes, dislikes, err
}

func GetLikesForComment(db *sql.DB, commentID int) (int, int, error) {
	var likes, dislikes int
	err := db.QueryRow("SELECT COUNT(CASE WHEN like = 1 THEN 1 END), COUNT(CASE WHEN like = 0 THEN 1 END) FROM likes WHERE comment_id = ?", commentID).Scan(&likes, &dislikes)
	return likes, dislikes, err
}
