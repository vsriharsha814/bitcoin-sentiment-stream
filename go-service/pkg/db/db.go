package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"
	"fmt"
	"strings"

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

// FetchInitialLastSentiments returns the most recent sentiment_score
// for each coin in coinIDs, before the given time.
// It uses a single DISTINCT ON query.
func FetchInitialLastSentiments(coinIDs []int, before time.Time) (map[int]float64, error) {
    // build a SQL placeholder list: ($1,$2, …)
    placeholders := make([]string, len(coinIDs))
    args := make([]interface{}, len(coinIDs)+1)
    for i, id := range coinIDs {
        placeholders[i] = fmt.Sprintf("$%d", i+1)
        args[i] = id
    }
    // the last arg is the `before` timestamp
    args[len(coinIDs)] = before

    sql := fmt.Sprintf(`
        SELECT DISTINCT ON (coin_id)
               coin_id, sentiment_score
          FROM aggregated_sentiments
         WHERE coin_id IN (%s)
           AND window_start < $%d
         ORDER BY coin_id, window_start DESC
    `, strings.Join(placeholders, ","), len(coinIDs)+1)

    rows, err := Conn.QueryContext(context.Background(), sql, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    result := make(map[int]float64, len(coinIDs))
    for rows.Next() {
        var cid int
        var score float64
        if err := rows.Scan(&cid, &score); err != nil {
            return nil, err
        }
        result[cid] = score
    }
    return result, rows.Err()
}

// InsertAggregatedSentimentBatch bulk-inserts all new records at once.
func InsertAggregatedSentimentBatch(records []struct {
    CoinID     int
    Window     time.Time
    Sentiment  float64
}) error {
    if len(records) == 0 {
        return nil
    }
    // build a VALUES list: ($1,$2,$3),($4,$5,$6),…
    var placeholders []string
    args := make([]interface{}, 0, len(records)*3)
    for i, rec := range records {
        base := i*3 + 1
        placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d)",
            base, base+1, base+2))
        args = append(args, rec.CoinID, rec.Window, rec.Sentiment)
    }
    sql := fmt.Sprintf(`
        INSERT INTO aggregated_sentiments
          (coin_id, window_start, sentiment_score)
        VALUES %s
        ON CONFLICT (coin_id, window_start) DO NOTHING
    `, strings.Join(placeholders, ","))

    _, err := Conn.ExecContext(context.Background(), sql, args...)
    return err
}

// FetchRawMessagesBetween returns every raw_messages row whose created_at
// is ≥ start AND < end, ordered oldest→newest.
func FetchRawMessagesBetween(start, end time.Time) ([]MessageScore, error) {
    const q = `
      SELECT
          question_id,
          currency_id,
          sentiment_score,
          created_at
        FROM raw_messages
       WHERE created_at >= $1
         AND created_at <  $2
       ORDER BY created_at ASC
    `
    rows, err := Conn.QueryContext(context.Background(), q, start, end)
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
    return out, rows.Err()
}

// FetchLastAggregatedSentiment returns the most recent sentiment_score
// for a single coin before the given time (or 0 if none).
func FetchLastAggregatedSentiment(coinID int, before time.Time) (float64, error) {
    // reuse the bulk helper for a single-element slice
    m, err := FetchInitialLastSentiments([]int{coinID}, before)
    if err != nil {
        return 0, err
    }
    if v, ok := m[coinID]; ok {
        return v, nil
    }
    return 0, nil
}

// AggregatedSentiment holds one row from aggregated_sentiments.
type AggregatedSentiment struct {
    CurrencyID     int
    WindowStart    time.Time
    SentimentScore float64
}

// FetchAggregatedSentimentsBetween returns all aggregated_sentiments
// rows whose window_start is in [start, end], oldest first.
func FetchAggregatedSentimentsBetween(start, end time.Time) ([]AggregatedSentiment, error) {
    const q = `
      SELECT coin_id, window_start, sentiment_score
        FROM aggregated_sentiments
       WHERE window_start >= $1
         AND window_start <= $2
       ORDER BY window_start ASC
    `
    rows, err := Conn.QueryContext(context.Background(), q, start, end)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []AggregatedSentiment
    for rows.Next() {
        var a AggregatedSentiment
        if err := rows.Scan(&a.CurrencyID, &a.WindowStart, &a.SentimentScore); err != nil {
            return nil, err
        }
        out = append(out, a)
    }
    return out, rows.Err()
}

