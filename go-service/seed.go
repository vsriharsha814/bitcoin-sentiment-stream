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

	// Create questions table
	createQuestions := `
	CREATE TABLE IF NOT EXISTS questions (
	  id   SERIAL PRIMARY KEY,
	  code VARCHAR(4) NOT NULL UNIQUE,
	  text TEXT NOT NULL
	);
	`
	if _, err := db.ExecContext(ctx, createQuestions); err != nil {
		log.Fatalf("Failed to create questions table: %v", err)
	}
	log.Println("Ensured questions table exists")

	// Insert question mappings
	mapping := map[string]string{
		"1": "New Features or Use Cases of \"coin_name\"",
		"2": "Founders or Leadership of \"coin_name\"",
		"3": "Security Concerns or Hacks related to \"coin_name\"",
		"4": "Market Trends and Price Predictions of \"coin_name\"",
		"5": "Regulatory Updates and Government Policies affecting \"coin_name\"",
		"6": "Community Sentiment and Adoption for \"coin_name\"",
		"7": "Partnerships and Integrations involving \"coin_name\"",
		"8": "Mining and Staking Discussions around \"coin_name\"",
	}
	insertQ := `
	INSERT INTO questions (code, text)
	VALUES ($1, $2)
	ON CONFLICT (code) DO NOTHING;
	`
	for code, text := range mapping {
		if _, err := db.ExecContext(ctx, insertQ, code, text); err != nil {
			log.Printf("Error inserting question mapping %s: %v", code, err)
		}
	}
	log.Println("Inserted question mappings")

	// Create raw_messages table with sentiment_score column
	createRaw := `
	CREATE TABLE IF NOT EXISTS raw_messages (
	  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	  source TEXT NOT NULL,
	  external_id TEXT NOT NULL UNIQUE,
	  question_id INTEGER NOT NULL REFERENCES questions(id),
	  currency_id INTEGER NOT NULL REFERENCES currency(id),
	  author TEXT,
	  content TEXT NOT NULL,
	  sentiment_score FLOAT,  -- Added sentiment_score column
	  created_at TIMESTAMPTZ NOT NULL,
	  fetched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	  metadata JSONB DEFAULT '{}'::jsonb
	);
	`
	if _, err := db.ExecContext(ctx, createRaw); err != nil {
		log.Fatalf("Failed to create raw_messages table: %v", err)
	}
	log.Println("Ensured raw_messages table exists")

	// Sample data to insert (with sentiment scores)
	samples := []struct {
		Source       string
		ExternalID   string
		QuestionCode string
		CurrencyID   int
		Author       string
		Content      string
		Sentiment    float64
		CreatedAt    time.Time
	}{
		{"twitter", "tweet1", "1", 91, "@alice", "Check out new features of coin_name!", 0.75, time.Now().Add(-10 * time.Minute)},
		{"reddit",  "post1", "2", 92, "u/bob", "Leadership changes announced for coin_name.", 0.65, time.Now().Add(-5 * time.Minute)},
		{"twitter", "tweet2", "3", 93, "@john", "Big concerns over security flaws in coin_name.", -0.85, time.Now().Add(-15 * time.Minute)},
		{"reddit",  "post2", "4", 94, "u/susan", "Market trends suggest a bullish run for coin_name.", 0.80, time.Now().Add(-20 * time.Minute)},
		{"twitter", "tweet3", "5", 95, "@dave", "Regulatory updates on coin_name show some new compliance requirements.", 0.55, time.Now().Add(-25 * time.Minute)},
		{"reddit",  "post3", "6", 96, "u/paul", "The community is really excited about coin_name’s upcoming feature!", 0.90, time.Now().Add(-30 * time.Minute)},
		{"twitter", "tweet4", "7", 97, "@jane", "coin_name has announced new partnerships with major platforms.", 0.70, time.Now().Add(-35 * time.Minute)},
		{"reddit",  "post4", "8", 98, "u/ted", "Discussions about staking rewards for coin_name are heating up.", 0.60, time.Now().Add(-40 * time.Minute)},
	}

	// Prepare insert statement for raw_messages
	insertRaw := `
	INSERT INTO raw_messages
	  (source, external_id, question_id, currency_id, author, content, sentiment_score, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (external_id) DO NOTHING;
	`

	for _, s := range samples {
		// lookup question_id
		var qid int
		if err := db.QueryRowContext(ctx, "SELECT id FROM questions WHERE code=$1", s.QuestionCode).Scan(&qid); err != nil {
			log.Printf("Unknown question code %s: %v", s.QuestionCode, err)
			continue
		}

		if _, err := db.ExecContext(ctx, insertRaw,
			s.Source, s.ExternalID, qid, s.CurrencyID,
			s.Author, s.Content, s.Sentiment, s.CreatedAt); err != nil {
			log.Printf("Error inserting raw message (external_id=%s): %v", s.ExternalID, err)
		} else {
			fmt.Printf("Inserted raw message: external_id=%s\n", s.ExternalID)
		}
	}
}
