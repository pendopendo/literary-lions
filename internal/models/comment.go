package models

//(27.09.24)
import (
	"database/sql"
	//"errors" // New import
)

// Define a Comment type to hold the data for an individual Comment. Notice how
// the fields of the struct correspond to the fields in our MySQL Comments
// table?
type Comment struct {
	ID       int
	User     User
	Post     int
	Text     string
	Created  string
	Likes    int
	Dislikes int
}

// Define a CommentModel type which wraps a sql.DB connection pool.
type CommentModel struct {
	DB *sql.DB
}

//funtkstioonid

func (m *CommentModel) GetCommentsForPost(postid int) ([]Comment, error) {
	// Write the SQL statement we want to execute.
	statement := `SELECT id, user, post, text, created FROM comment WHERE post=? ORDER BY id DESC`

	// Use the Query() method on the connection pool to execute our SQL statement.
	rows, err := m.DB.Query(statement, postid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize an empty slice to hold the Comment structs.
	var comments []Comment

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		var s Comment
		var userID int

		// Scan the row data into our Comment struct.
		err = rows.Scan(&s.ID, &userID, &s.Post, &s.Text, &s.Created)
		if err != nil {
			return nil, err
		}

		// Fetch the User details.
		users := &UserModel{DB: m.DB}
		s.User, err = users.Get(userID)
		if err != nil {
			return nil, err
		}

		// Fetch the likes and dislikes for the comment.
		err = m.getLikesDislikes(&s)
		if err != nil {
			return nil, err
		}

		// Append the comment to our slice.
		comments = append(comments, s)
	}

	// Check for any errors encountered during iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Return the slice of comments.
	return comments, nil
}

// Helper method to fetch likes and dislikes for a comment.
func (m *CommentModel) getLikesDislikes(comment *Comment) error {
	// Query to count the likes for the comment.
	likeStmt := `SELECT COUNT(*) FROM reactions WHERE comment=? AND reaction='like'`
	err := m.DB.QueryRow(likeStmt, comment.ID).Scan(&comment.Likes)
	if err != nil {
		return err
	}

	// Query to count the dislikes for the comment.
	dislikeStmt := `SELECT COUNT(*) FROM reactions WHERE comment=? AND reaction='dislike'`
	err = m.DB.QueryRow(dislikeStmt, comment.ID).Scan(&comment.Dislikes)
	if err != nil {
		return err
	}

	return nil
}

// proovin lisada commentit vist siin
func (m *CommentModel) InsertComment(postiID int, text string, userID int) error {
	// Write the SQL statement we want to execute.

	statement := `
	INSERT INTO comment (user, post, text, created) 
	VALUES (?, ?, ?, datetime('now', '+3 hours'));
	`

	// Execute the statement
	_, err := m.DB.Exec(statement, userID, postiID, text)
	if err != nil {
		return err
	}

	// // Get the ID of the newly inserted row
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	// Return the ID as an int
	return nil
}
