package database

// TODO: Replace SQLite with PostgreSQL and integrate a migration tool (e.g., golang-migrate)
// This InitDB() is a temporary setup for MVP. Migrations will allow safer schema changes later.
// This file initializes the database connection and creates necessary tables.
import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQlite3 driver
	//_ "github.com/lib/pq" // PostgreSQL driver, if needed
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error
	// Open a new database connection
	DB, err = sql.Open("sqlite3", "echoes.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Check if the database is reachable
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create the users table if it doesn't exist
	// This is a simple example, will use migrations in a real application
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);
	`

	_, err = DB.Exec(createUserTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
	// Log successful connection
	log.Println("Database connection established successfully")
	fmt.Println("Database connection established successfully")
}
