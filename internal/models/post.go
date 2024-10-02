package models

//(27.09.24)
import (
	"database/sql"
	//"errors" // New import
	// "time"
	"errors"
)

// Define a Post type to hold the data for an individual Post. Notice how
// the fields of the struct correspond to the fields in our MySQL Posts
// table?
type Post struct {
	ID       int
	User     User
	Title    string
	Text     string
	Category int
	Created  string
}

// Define a PostModel type which wraps a sql.DB connection pool.
type PostModel struct {
	DB *sql.DB
}

//funtkstioonid

// kõik postitused saab
func (m *PostModel) Get() ([]Post, error) {
	// Write the SQL statement we want to execute.
	statement := `SELECT id, user, title, text, category, created FROM post ORDER BY id DESC`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	rows, err := m.DB.Query(statement)

	//fmt.Println(rows)

	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the Post structs.
	var posts []Post

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Post struct.
		var s Post
		var userID int
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Post object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &userID, &s.Title, &s.Text, &s.Category, &s.Created) //igalt realt loeb ja panne slice
		if err != nil {
			return nil, err
		}

		users := &UserModel{DB: m.DB}
		s.User, err = users.Get(userID)
		if err != nil {
			return nil, err
		}

		// Append it to the slice of posts.
		posts = append(posts, s) //siin paneb slice
		//fmt.Println("rida:", posts)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Posts slice.
	return posts, nil
}

func (m *PostModel) GetRow(id int) (Post, error) {
	// Write the SQL statement we want to execute.
	statement := `SELECT id, user, title, text, category, created FROM post where ID=?`
	// siin on placeholder parametr
	//päring andmebaasi ÜHELE reale
	row := m.DB.QueryRow(statement, id)

	var s Post
	var userID int

	err := row.Scan(&s.ID, &userID, &s.Title, &s.Text, &s.Category, &s.Created)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return Post{}, ErrNoRecord
		} else {

			return Post{}, err
		}
	}

	users := &UserModel{DB: m.DB}
	s.User, err = users.Get(userID)
	if err != nil {
		return Post{}, err
	}

	return s, nil

}

//atahame saada kategooria ID järgi sellele vastavad postitused

func (m *PostModel) GetPostsForCategory(postid int) ([]Post, error) {
	// Write the SQL statement we want to execute.
	statement := `SELECT id, user, title, text, created FROM post where category=? ORDER BY id DESC`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	rows, err := m.DB.Query(statement, postid) //siit ülealt. Query annab mingisuguse vastuse

	//fmt.Println(rows)

	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the Post structs.
	var posts []Post

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Post struct.
		var s Post
		var userID int
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Post object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		//id, user, post, text, created
		err = rows.Scan(&s.ID, &userID, &s.Title, &s.Text, &s.Created) //igalt realt loeb ja panne slice
		if err != nil {
			return nil, err
		}
		users := &UserModel{DB: m.DB}
		s.User, err = users.Get(userID)
		if err != nil {
			return nil, err
		}

		// Append it to the slice of posts.
		posts = append(posts, s) //siin paneb slice
		//fmt.Println("rida:", posts)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Posts slice.
	return posts, nil
}

func (m *PostModel) Insert(categoryID int, title string, text string, userID int) (int, error) {
	// Write the SQL statement we want to execute.
	// siin all on
	// (1, ?, ?, 1, '2024-10-01 10:30:00');
	// us siis 1 - user id, ? title, ? text (neid tahame lisada), 1 - kategooria id, created - hetkel siin pandud
	//SELECT datetime('now'); - peaks võtma praeguse aja
	statement := `
	INSERT INTO post (user, title, text, category, created) 
	VALUES (?, ?, ?, ?, datetime('now', '+3 hours'));
	`

	// Execute the statement
	result, err := m.DB.Exec(statement, userID, title, text, categoryID)
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
