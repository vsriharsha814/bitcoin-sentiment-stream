package handlers

import (
	"log"
	"net/http"
	"time"
	"strconv"
	"sort"

	"github.com/gorilla/websocket"

	
	// "github.com/cosmic-hash/CryptoPulse/pkg/config"
	// "github.com/cosmic-hash/CryptoPulse/pkg/model"
	"github.com/cosmic-hash/CryptoPulse/pkg/db"

)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSHandler streams sentiment grouped by UTC minute.
func WSHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade the HTTP connection to a WebSocket.
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }
    defer conn.Close()
    log.Println("WebSocket connection established")

    send := func() {
        // 1) Fetch all recent rows from DB.
        samples, err := db.FetchMessageScoresFromDB()
        if err != nil {
            log.Printf("DB fetch error, skipping update: %v", err)
            return
        }
        log.Printf("Fetched %d rows: %+v", len(samples), samples)

        // 2) Group samples into buckets by UTC minute.
        buckets := make(map[time.Time][]db.MessageScore)
        for _, m := range samples {
            // Truncate to the minute in UTC.
            bucketTime := m.CreatedAt.UTC().Truncate(time.Minute)
            buckets[bucketTime] = append(buckets[bucketTime], m)
        }
        log.Printf("Buckets grouped by minute: %+v", buckets)

        // 3) Sort bucket times ascending.
        times := make([]time.Time, 0, len(buckets))
        for t := range buckets {
            times = append(times, t)
        }
        sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })
        log.Printf("Sorted bucket times: %+v", times)

        // 4) Build the response array.
        resp := make([]map[string]interface{}, 0, len(times))
        for _, t := range times {
            group := buckets[t]
            // Sum & count per currency.
            sums := map[int]float64{}
            counts := map[int]int{}
            for _, m := range group {
                sums[m.CurrencyID] += m.SentimentScore
                counts[m.CurrencyID]++
            }
            log.Printf("Window %s sums: %+v counts: %+v", t.Format(time.RFC3339), sums, counts)

            // Build coins array for this bucket.
            coins := make([]map[string]float64, 0, len(sums))
            for cid, sum := range sums {
                avg := sum / float64(counts[cid])
                coins = append(coins, map[string]float64{
                    strconv.Itoa(cid): avg,
                })
            }
            log.Printf("Coins for %s: %+v", t.Format(time.RFC3339), coins)

            // Format time as "YYYY-MM-DDTHH:MMZ" (UTC).
            entry := map[string]interface{}{
                "time":  t.Format("2006-01-02T15:04Z"),
                "coins": coins,
            }
            resp = append(resp, entry)
        }
        log.Printf("Final response payload: %+v", resp)

        // 5) Send the JSON payload over WebSocket.
        if err := conn.WriteJSON(resp); err != nil {
            log.Println("WebSocket write error:", err)
        } else {
            log.Println("WebSocket: sent update")
        }
    }

    // Initial send.
    send()
    // Then every 5 minutes.
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    for range ticker.C {
        send()
    }
}