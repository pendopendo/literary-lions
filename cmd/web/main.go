package main

/*
http://localhost:4000/

http://localhost:4000/post/view/andre

kontrollid

http://localhost:4000/post/view/1
http://localhost:4000/post/create
http://localhost:4000/comment/create
http://localhost:4000/user/signup
http://localhost:4000/user/login

*/
import (
	"database/sql" // New import
	//"flag"
	"log"
	"net/http"
	"os"

	"log/slog"
	"time" // New import

	"snippetbox.alexedwards.net/internal/models"

	"github.com/alexedwards/scs/sqlite3store" // Import the SQLite3 store
	"github.com/alexedwards/scs/v2"           // The SCS session manager package
	_ "github.com/mattn/go-sqlite3"           // SQLite3 driver
	// "fmt"
)

// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
// kui vastavas go failis  loodud  nt Get funktsioon siis siin peab selle lisama
// ja ka siia alla app := application lisada
type application struct {
	logger         *slog.Logger
	categories     *models.CategoryModel
	posts          *models.PostModel
	comments       *models.CommentModel //sID 240928_184230
	users          *models.UserModel
	sessionManager *scs.SessionManager
	reactions      *models.ReactionModel
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Successfully connected to the SQLite database!")

	// Connect to the database using the connectDB function
	db, err := connectDB("./literarylionforum.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Use the scs.New() function to initialize a new session manager. Then we
	// configure it to use our MySQL database as the session store, and set a
	// lifetime of 12 hours (so that sessions automatically expire 12 hours
	// after first being created).
	sessionManager := scs.New()
	// Assuming 'db' is your SQLite3 database connection
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 1 * time.Hour //5 mintisiis lgib ise välhja

	// Initialize a models.SnippetModel instance containing the connection pool
	// and add it to the application dependencies.
	//application asju saavad kõik osad kasutada
	app := &application{
		logger:         logger,
		categories:     &models.CategoryModel{DB: db}, //db ülevalt
		posts:          &models.PostModel{DB: db},     //db ülevalt
		comments:       &models.CommentModel{DB: db},  //db ülevalt
		users:          &models.UserModel{DB: db},
		sessionManager: sessionManager,
		reactions:      &models.ReactionModel{DB: db},
	}

	logger.Info("starting server on :8080")

	err = http.ListenAndServe(":8080", app.routes())
	log.Fatal(err)

}

// connectDB opens a connection to the SQLite database and returns the database instance
func connectDB(dbPath string) (*sql.DB, error) {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err = db.Ping(); err != nil {
		db.Close() // Close the connection if Ping fails
		return nil, err
	}

	return db, nil
}

// see korraks
// getPost retrieves a post by its ID
func getPost(db *sql.DB, postID int) (int, string, string, int, string, error) {
	var id int
	var title, text, created string
	var userID, categoryID int

	query := `SELECT id, user, title, text, category, created FROM post WHERE id = ?`
	err := db.QueryRow(query, postID).Scan(&id, &userID, &title, &text, &categoryID, &created)
	if err != nil {
		return 0, "", "", 0, "", err
	}

	return id, title, text, categoryID, created, nil
}

// ------------------------------------ //
