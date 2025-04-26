package handlers

import (
    "log"
    "net/http"
    "sort"
    "strconv"
    "strings"
    "time"

    "github.com/gorilla/websocket"
    "github.com/cosmic-hash/CryptoPulse/pkg/db"
)

// upgrader allows HTTP → WebSocket upgrade
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

// build a fast lookup from your global coinsList in api.go
var currencyCodeMap = func() map[int]string {
    m := make(map[int]string, len(coinsList))
    for _, c := range coinsList {
        m[c.ID] = c.Code
    }
    return m
}()

// WSHandler streams pre-aggregated sentiment in 5-minute buckets.
// Supports both query-param and JSON overrides.
// If no "tokens" key is sent, it will send all coins.
func WSHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("[WS] upgrade error:", err)
        return
    }
    defer conn.Close()
    log.Println("[WS] connection established")

    // --- initial overrides from query params ---
    var (
        filterCodes []string
        useFilter   bool
        fixedStart  time.Time
        fixedEnd    time.Time
        useFixed    bool
    )
    if tok := r.URL.Query().Get("tokens"); tok != "" {
        filterCodes = strings.Split(tok, ",")
        useFilter = true
        log.Printf("[WS] initial token filter: %v", filterCodes)
    }
    if s := r.URL.Query().Get("start_time"); s != "" {
        if t, err := time.Parse(time.RFC3339, s); err == nil {
            fixedStart = t.UTC()
            useFixed = true
            log.Printf("[WS] initial start_time override: %s", fixedStart)
        }
    }
    if e := r.URL.Query().Get("end_time"); e != "" {
        if t, err := time.Parse(time.RFC3339, e); err == nil {
            fixedEnd = t.UTC()
            useFixed = true
            log.Printf("[WS] initial end_time override: %s", fixedEnd)
        }
    }

    // channel to trigger immediate refresh
    overrideCh := make(chan struct{}, 1)

    // --- reader goroutine for JSON control frames ---
    go func() {
        for {
            var msg struct {
                Tokens    *[]string `json:"tokens"`
                StartTime *string   `json:"start_time"`
                EndTime   *string   `json:"end_time"`
            }
            if err := conn.ReadJSON(&msg); err != nil {
                log.Println("[WS] read JSON error:", err)
                close(overrideCh)
                return
            }
            log.Printf("[WS] got override frame: %+v", msg)

            // tokens logic: nil = no key, empty slice = explicit empty
            if msg.Tokens != nil {
                filterCodes = *msg.Tokens
                useFilter = true
                log.Printf("[WS] updated token filter: %v", filterCodes)
            } else {
                useFilter = false
                log.Println("[WS] no tokens key → sending ALL coins")
            }

            // start_time override
            if msg.StartTime != nil {
                if t, err := time.Parse(time.RFC3339, *msg.StartTime); err == nil {
                    fixedStart = t.UTC()
                    useFixed = true
                    log.Printf("[WS] updated start_time: %s", fixedStart)
                } else {
                    log.Printf("[WS] bad start_time %q: %v", *msg.StartTime, err)
                }
            }
            // end_time override
            if msg.EndTime != nil {
                if t, err := time.Parse(time.RFC3339, *msg.EndTime); err == nil {
                    fixedEnd = t.UTC()
                    useFixed = true
                    log.Printf("[WS] updated end_time: %s", fixedEnd)
                } else {
                    log.Printf("[WS] bad end_time %q: %v", *msg.EndTime, err)
                }
            }

            // fire an immediate refresh
            select {
            case overrideCh <- struct{}{}:
            default:
            }
        }
    }()

    // helper: generate every 5-min tick between start and end
    makeTimeline := func(start, end time.Time) []time.Time {
        start = start.Truncate(5 * time.Minute)
        var series []time.Time
        for t := start; !t.After(end); t = t.Add(5 * time.Minute) {
            series = append(series, t)
        }
        return series
    }

    // core send logic
    send := func() {
        // pick window
        var start, end time.Time
        if useFixed {
            start, end = fixedStart, fixedEnd
            log.Printf("[WS] using fixed window %s → %s", start, end)
        } else {
            end = time.Now().UTC()
            start = end.Add(-1 * time.Hour)
            log.Printf("[WS] using default window %s → %s", start, end)
        }

        // fetch aggregated rows
        aggs, err := db.FetchAggregatedSentimentsBetween(start, end)
        if err != nil {
            log.Printf("[WS] fetch error: %v", err)
            return
        }
        log.Printf("[WS] fetched %d rows", len(aggs))

        // bucket by minute → map[timestamp][code] = score
        buckets := make(map[time.Time]map[string]float64)
        for _, a := range aggs {
            ts := a.WindowStart.UTC().Truncate(time.Minute)
            if buckets[ts] == nil {
                buckets[ts] = make(map[string]float64)
            }
            code := strconv.Itoa(a.CurrencyID)
            if c, ok := currencyCodeMap[a.CurrencyID]; ok {
                code = c
            }
            buckets[ts][code] = a.SentimentScore
        }

        // build full timeline
        timeline := makeTimeline(start, end)

        // determine which codes to include
        var codes []string
        if useFilter {
            codes = filterCodes
        } else {
            for _, c := range coinsList {
                codes = append(codes, c.Code)
            }
        }
        sort.Strings(codes)

        // assemble payload
        resp := make([]map[string]interface{}, 0, len(timeline))
        for _, ts := range timeline {
            data := make(map[string]float64, len(codes))
            bucket := buckets[ts]
            for _, code := range codes {
                data[code] = bucket[code] // zero if missing
            }
            log.Printf("[WS] bucket %s → %+v", ts.Format(time.RFC3339), data)
            resp = append(resp, map[string]interface{}{
                "time":  ts.Format("2006-01-02T15:04Z"),
                "coins": data,
            })
        }

        // send JSON
        if err := conn.WriteJSON(resp); err != nil {
            log.Println("[WS] write error:", err)
        } else {
            log.Printf("[WS] sent %d buckets", len(resp))
        }
    }

    // initial send, then every minute, plus immediate on override
    send()
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            send()
        case _, ok := <-overrideCh:
            if !ok {
                return
            }
            log.Println("[WS] override fired — immediate send")
            send()
        }
    }
}