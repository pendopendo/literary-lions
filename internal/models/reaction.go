package models

//(27.09.24)
import (
	"database/sql"
	//"errors" // New import
	"time"
)

// Define a Reaction type to hold the data for an individual Reaction. Notice how
// the fields of the struct correspond to the fields in our MySQL Reactions
// table?
type Reaction struct {
	ID       int
	Reaction string

	Comment int
	Post    int
	User    string
	Created time.Time
}

// Define a ReactionModel type which wraps a sql.DB connection pool.
type ReactionModel struct {
	DB *sql.DB
}

// funtkstioonid
// proovin lisada commentit vist siin
func (m *ReactionModel) CountPostLikesDislikes(postID int) (int, int, error) {
	// Write the SQL statement to count likes and dislikes for the given postID
	statement := `
    SELECT 
        COALESCE(SUM(CASE WHEN reaction = 'like' THEN 1 ELSE 0 END), 0) AS likes,	
        COALESCE(SUM(CASE WHEN reaction = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes
    FROM reactions
    WHERE post = ?;
    `

	// Variables to store the count of likes and dislikes
	var likes, dislikes int

	// Execute the query and scan the results into likes and dislikes
	err := m.DB.QueryRow(statement, postID).Scan(&likes, &dislikes)
	if err != nil {
		return 0, 0, err
	}

	// Return the counts of likes and dislikes
	return likes, dislikes, nil
}

func (m *ReactionModel) InsertReaction(reaction string, postID int, userID int) (int, error) {
	// Write the SQL statement to check if the user has already reacted to the post.
	existingStatement := `
	SELECT id
	FROM reactions
	WHERE post = ? AND user = ?;
	`

	var existingID int
	err := m.DB.QueryRow(existingStatement, postID, userID).Scan(&existingID)

	// If error is different from no rows error, return the error
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	// If err is nil, it means the user has already reacted so we should update the reaction
	if existingID > 0 {
		updateStatement := `
		UPDATE reactions
		SET reaction = ?
		WHERE id = ?;
		`
		_, err := m.DB.Exec(updateStatement, reaction, existingID)
		if err != nil {
			return 0, err
		}
		return existingID, nil
	}

	// If no rows are returned, we can proceed with the insert
	insertStatement := `
	INSERT INTO reactions (reaction, post, user) 
	VALUES (?, ?, ?);
	`

	// Execute the insert statement
	result, err := m.DB.Exec(insertStatement, reaction, postID, userID)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted row
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Return the ID as an int
	return int(id), nil
}

func (m *ReactionModel) InsertReactionComment(reaction string, commentID int, userID int) (int, error) {
	// Write the SQL statement to check if the user has already reacted to the post.
	existingStatement := `
	SELECT id
	FROM reactions
	WHERE comment = ? AND user = ?;
	`

	var existingID int
	err := m.DB.QueryRow(existingStatement, commentID, userID).Scan(&existingID)

	// If error is different from no rows error, return the error
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	// If err is nil, it means the user has already reacted so we should update the reaction
	if existingID > 0 {
		updateStatement := `
		UPDATE reactions
		SET reaction = ?
		WHERE id = ?;
		`
		_, err := m.DB.Exec(updateStatement, reaction, existingID)
		if err != nil {
			return 0, err
		}
		return existingID, nil
	}

	// If no rows are returned, we can proceed with the insert
	insertStatement := `
	INSERT INTO reactions (reaction, comment, user) 
	VALUES (?, ?, ?);
	`

	// Execute the insert statement
	result, err := m.DB.Exec(insertStatement, reaction, commentID, userID)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted row
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Return the ID as an int
	return int(id), nil
}
