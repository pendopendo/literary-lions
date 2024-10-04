package models

//(27.09.24)
import (
	"database/sql"
	//"errors" // New import
	"fmt"
)

// Define a category type to hold the data for an individual category. Notice how
// the fields of the struct correspond to the fields in our MySQL categorys
// table?
type Category struct {
	ID   int
	Name string
}

// Define a categoryModel type which wraps a sql.DB connection pool.
type CategoryModel struct {
	DB *sql.DB
}

//funtkstioonid

// This will return a specific Category based on its id.
func (m *CategoryModel) GetRow() ([]Category, error) {
	//statement
	statement := "select ID, name from category" //kõik kategooriad
	// siin on placeholder parametr
	//päring andmebaasi ÜHELE reale
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

	// Initialize an empty slice to hold the Category structs.
	var categories []Category

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Category struct.
		var s Category
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Category object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		//id, user, post, text, created
		err = rows.Scan(&s.ID, &s.Name) //igalt realt loeb ja panne slice
		if err != nil {
			return nil, err
		}
		// Append it to the slice of comments.
		categories = append(categories, s) //siin paneb slice
		fmt.Println("rida kat:", categories)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Categorys slice.
	return categories, nil

	
}
// GetByID returns a specific category based on its ID.
func (m *CategoryModel) GetByID(id int) (*Category, error) {
	stmt := `SELECT id, name FROM category WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)

	// Loome uue kategooria struktuuri ja täidame selle andmetega
	category := &Category{}
	err := row.Scan(&category.ID, &category.Name)
	if err == sql.ErrNoRows {
		return nil, ErrNoRecord // Saadad tagasi vea, kui kirjet ei leitud
	} else if err != nil {
		return nil, err
	}

	return category, nil
}