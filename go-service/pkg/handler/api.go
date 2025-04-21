package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	// "github.com/cosmic-hash/CryptoPulse/pkg/config"
"github.com/cosmic-hash/CryptoPulse/pkg/db"
	// "github.com/cosmic-hash/CryptoPulse/pkg/model"
)

// HelloHandler serves GET /
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Hello"))
}

// SentimentHandler serves POST /sentiment
func SentimentHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // 1) Fetch data from DB
    samples, err := db.FetchMessageScoresFromDB()
    if err != nil {
        log.Printf("DB fetch error: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    log.Printf("Fetched %d samples", len(samples))

    // 2) Prepare three 5‑minute buckets ending now, now‑5min, now‑10min
    now := time.Now()
    bucketEnds := []time.Time{
        now.Add(-10 * time.Minute),
        now.Add(-5 * time.Minute),
        now,
    }

    // 3) Build response array
    response := make([]map[string]interface{}, 0, len(bucketEnds))
    for _, end := range bucketEnds {
        start := end.Add(-5 * time.Minute)
        entry := map[string]interface{}{
            "time": end.Format("02-January-2006 15:04"),
        }

        // sum & count per currency
        sums := map[int]float64{}
        counts := map[int]int{}
        for _, m := range samples {
            if m.CreatedAt.After(start) && !m.CreatedAt.After(end) {
                sums[m.CurrencyID] += m.SentimentScore
                counts[m.CurrencyID]++
            }
        }

        // compute averages
        for cid, sum := range sums {
            avg := sum / float64(counts[cid])
            entry[strconv.Itoa(cid)] = avg
        }
        response = append(response, entry)
    }

    // 4) Return JSON
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Encode error: %v", err)
    }
}