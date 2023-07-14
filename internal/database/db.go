package mydatabase

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"service1/internal/models"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func Connect() (*DB, error) {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database")
	return &DB{db}, nil
}

func (db *DB) CreateTableIfNotExists() error {
	query := `
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			body TEXT,
			user_id INT
			-- Add other necessary columns here
		)
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Table created successfully")
	return nil
}

func (db *DB) SavePost(post *models.Post) error {
	query := "INSERT INTO posts (id, user_id, title, body) VALUES ($1, $2, $3, $4)"

	_, err := db.Exec(query, post.ID, post.UserID, post.Title, post.Body)
	if err != nil {
		return err
	}

	log.Printf("Saved post with ID: %d", post.ID)
	return nil
}
