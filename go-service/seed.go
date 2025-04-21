package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Load environment variables from .env file, if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Read DATABASE_URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Connect to Postgres using pgx driver
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Verify connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("DB unreachable: %v", err)
	}

	// Create table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS raw_messages (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	source TEXT NOT NULL,
	external_id TEXT NOT NULL UNIQUE,
	question_code VARCHAR(4),
	author TEXT,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	metadata JSONB DEFAULT '{}'::jsonb
);`

	if _, err := db.ExecContext(ctx, createTableSQL); err != nil {
		log.Fatalf("Failed to create raw_messages table: %v", err)
	}
	log.Println("Ensured raw_messages table exists")

	// Sample data to insert
	samples := []struct {
		Source       string
		ExternalID   string
		QuestionCode string
		Author       string
		Content      string
		CreatedAt    time.Time
	}{
		{"twitter", "tweet1", "1", "@alice", "Check out new features of coin_name!", time.Now().Add(-10 * time.Minute)},
		{"reddit",  "post1",  "2", "u/bob",  "Leadership changes announced for coin_name.", time.Now().Add(-5 * time.Minute)},
		// Add more entries as needed
	}

	// Prepare INSERT statement
	insertSQL := `INSERT INTO raw_messages
	(source, external_id, question_code, author, content, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (external_id) DO NOTHING;`

	// Execute inserts
	for _, s := range samples {
		if _, err := db.ExecContext(ctx, insertSQL,
			s.Source, s.ExternalID, s.QuestionCode, s.Author, s.Content, s.CreatedAt); err != nil {
			log.Printf("Error inserting sample (source=%s, external_id=%s): %v", s.Source, s.ExternalID, err)
		} else {
			fmt.Printf("Inserted sample: source=%s, external_id=%s\n", s.Source, s.ExternalID)
		}
	}
}
