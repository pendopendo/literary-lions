package models

import (
    "database/sql"
    "time"
)

type Comment struct {
    ID        int
    Content   string
    UserID    int
    PostID    int
    CreatedAt time.Time
}

func CreateComment(db *sql.DB, content string, userID, postID int) error {
    _, err := db.Exec("INSERT INTO comments (content, user_id, post_id) VALUES (?, ?, ?)", content, userID, postID)
    return err
}

func GetCommentsByPost(db *sql.DB, postID int) ([]Comment, error) {
    rows, err := db.Query("SELECT id, content, user_id, post_id, created_at FROM comments WHERE post_id = ?", postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    comments := []Comment{}
    for rows.Next() {
        var comment Comment
        err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt)
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }

    return comments, nil
}
