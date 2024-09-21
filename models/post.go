package models

import (
    "database/sql"
    "time"
)

type Post struct {
    ID        int
    Title     string
    Content   string
    UserID    int
    CategoryID int
    CreatedAt time.Time
}

func CreatePost(db *sql.DB, title, content string, userID, categoryID int) error {
    _, err := db.Exec("INSERT INTO posts (title, content, user_id, category_id) VALUES (?, ?, ?, ?)", title, content, userID, categoryID)
    return err
}

func GetPostsByCategory(db *sql.DB, categoryID int) ([]Post, error) {
    rows, err := db.Query("SELECT id, title, content, user_id, category_id, created_at FROM posts WHERE category_id = ?", categoryID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID, &post.CreatedAt)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}

func GetPostsByUser(db *sql.DB, userID int) ([]Post, error) {
    rows, err := db.Query("SELECT id, user_id, title, content, category_id, created_at FROM posts WHERE user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CategoryID, &post.CreatedAt); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func GetLikedPosts(db *sql.DB, userID int) ([]Post, error) {
    rows, err := db.Query(`
        SELECT p.id, p.user_id, p.title, p.content, p.category_id, p.created_at
        FROM posts p
        INNER JOIN likes l ON p.id = l.post_id
        WHERE l.user_id = ? AND l.is_like = true`, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CategoryID, &post.CreatedAt); err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }
    return posts, nil
}

func GetAllPosts(db *sql.DB) ([]Post, error) {
    rows, err := db.Query("SELECT id, title, content, user_id, category_id, created_at FROM posts")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    posts := []Post{}
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID, &post.CreatedAt)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}

