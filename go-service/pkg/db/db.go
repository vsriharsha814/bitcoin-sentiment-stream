package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var Conn *sql.DB

// InitDB initializes the global Conn handle.
func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}
	var err error
	Conn, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	Conn.SetMaxOpenConns(10)
	Conn.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Conn.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL")
}

// MessageScore holds a single row.
type MessageScore struct {
    QuestionID     string
    CurrencyID     int
    SentimentScore float64
    CreatedAt      time.Time
}

// FetchMessageScoresFromDB pulls question_id, currency_id, sentiment_score, created_at
// over the past 24 hours.
func FetchMessageScoresFromDB() ([]MessageScore, error) {
    const query = `
        SELECT
            question_id,
            currency_id,
            sentiment_score,
            created_at
          FROM raw_messages
         WHERE created_at >= NOW() - INTERVAL '24 HOUR'
         ORDER BY created_at DESC
    `
    rows, err := Conn.QueryContext(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []MessageScore
    for rows.Next() {
        var m MessageScore
        if err := rows.Scan(
            &m.QuestionID,
            &m.CurrencyID,
            &m.SentimentScore,
            &m.CreatedAt,
        ); err != nil {
            return nil, err
        }
        out = append(out, m)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    log.Printf("Fetched %d message scores\n", len(out))
    return out, nil
}