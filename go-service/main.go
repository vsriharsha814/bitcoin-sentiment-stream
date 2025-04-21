package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib" 
)

// global DB handle
var db *sql.DB

// Initialize the Postgres connection using the pgx driver.
func initDB() {
	dsn := os.Getenv("DATABASE_URL") // e.g. from Neon dashboard
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}
	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	// Optional tuning:
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Verify connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("✅ Connected to Neon PostgreSQL")
}

// fetchMessageScoresFromDB pulls the most recent sentiment scores grouped by question_id.
// Adjust the query to match your actual schema.
func fetchMessageScoresFromDB() (map[string][]float64, error) {
	query := `
		SELECT question_id, score
		FROM message_scores
		-- e.g. only last hour, or last 10 per question
		WHERE created_at >= NOW() - INTERVAL '1 HOUR'
		ORDER BY created_at DESC
	`
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scoresMap := make(map[string][]float64)
	for rows.Next() {
		var questionID string
		var score float64
		if err := rows.Scan(&questionID, &score); err != nil {
			return nil, err
		}
		scoresMap[questionID] = append(scoresMap[questionID], score)
	}
	return scoresMap, rows.Err()
}

// … your existing types and functions omitted for brevity …



// CalculateFinalSentiment computes the overall sentiment by aggregating average scores
// for each question multiplied by its corresponding weight.
func CalculateFinalSentiment(weights map[string]float64, messageScores map[string][]float64) float64 {
	finalSentiment := 0.0

	for question, weight := range weights {
		scores, exists := messageScores[question]
		if !exists || len(scores) == 0 {
			continue // Skip if no scores available for this question
		}
		sum := 0.0
		for _, score := range scores {
			sum += score
		}
		avgScore := sum / float64(len(scores))
		finalSentiment += weight * avgScore
	}
	return finalSentiment
}

// // sendSentimentUpdate now tries to fetch from DB first.
// func sendSentimentUpdate(conn *websocket.Conn) {
// 	// 1) Attempt to fetch from Neon
// 	numericScores, err := fetchMessageScoresFromDB()
// 	if err != nil {
// 		log.Printf("⚠️  DB fetch error, using defaults: %v", err)
// 		numericScores = defaultMessageScores
// 	}

// 	// 2) Map numeric keys → question names
// 	mappedScores := make(map[string][]float64)
// 	for key, scores := range numericScores {
// 		if qName, ok := questionMapping[key]; ok {
// 			mappedScores[qName] = scores
// 		}
// 	}

// 	// 3) Compute final sentiment
// 	finalSentiment := CalculateFinalSentiment(defaultWeights, mappedScores)
// 	resp := SentimentResponse{FinalSentiment: finalSentiment}

// 	// 4) Push over WebSocket
// 	if err := conn.WriteJSON(resp); err != nil {
// 		log.Println("WebSocket write error:", err)
// 	}
// }

// helloHandler is a basic handler that writes "Hello" to the response.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}

// SentimentRequest represents the JSON payload that carries sentiment scores.
// The keys are numbers (as strings) mapping to each question.
type SentimentRequest struct {
	MessageScores map[string][]float64 `json:"messageScores"`
}

// SentimentResponse represents the JSON response payload containing the computed sentiment.
type SentimentResponse struct {
	FinalSentiment float64 `json:"finalSentiment"`
}

// defaultWeights holds the weights for each question based on the provided questions.
// The keys here exactly match the question strings in the mapping file (with the placeholder).
var defaultWeights = map[string]float64{
	"New Features or Use Cases of \"coin_name\"":              0.15,
	"Founders or Leadership of \"coin_name\"":                 0.10,
	"Security Concerns or Hacks related to \"coin_name\"":     0.20,
	"Market Trends and Price Predictions of \"coin_name\"":    0.20,
	"Regulatory Updates and Government Policies affecting \"coin_name\"": 0.10,
	"Community Sentiment and Adoption for \"coin_name\"":      0.10,
	"Partnerships and Integrations involving \"coin_name\"":     0.10,
	"Mining and Staking Discussions around \"coin_name\"":      0.05,
}

// questionMapping holds the mapping from numeric keys (as strings) to question names.
// It will be loaded from mapping.json.
var questionMapping map[string]string

// defaultMessageScores is a simulated set of sentiment scores for each question.
// These use numeric keys and will be mapped to question text using questionMapping.
var defaultMessageScores = map[string][]float64{
	"1": {0.1, 0.2, 0.15, 0.1, 0.3, 0.05, 0.2, 0.15, 0.1, 0.2},
	"2": {0.2, 0.1, 0.15, 0.25, 0.1, 0.2, 0.15, 0.2, 0.1, 0.2},
	"3": {-0.1, -0.2, -0.15, -0.1, -0.3, -0.05, -0.2, -0.15, -0.1, -0.2},
	"4": {0.3, 0.25, 0.2, 0.35, 0.3, 0.25, 0.2, 0.35, 0.3, 0.25},
	"5": {0.0, 0.05, -0.05, 0.0, 0.1, 0.0, -0.1, 0.05, 0.0, 0.1},
	"6": {0.2, 0.3, 0.25, 0.2, 0.15, 0.2, 0.25, 0.3, 0.2, 0.15},
	"7": {0.1, 0.05, 0.1, 0.15, 0.1, 0.05, 0.1, 0.15, 0.1, 0.05},
	"8": {0.05, 0.0, 0.05, 0.1, 0.05, 0.0, 0.05, 0.1, 0.05, 0.0},
}

// sentimentHandler accepts a POST request with JSON input, converts numeric keys to question names,
// calculates the sentiment, and sends back the final sentiment as a JSON response.
func sentimentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SentimentRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Convert numeric keys in the payload to their corresponding question names using the mapping.
	mappedScores := make(map[string][]float64)
	for key, scores := range req.MessageScores {
		questionName, ok := questionMapping[key]
		if !ok {
			// Log unknown keys and skip them.
			log.Printf("Unknown question key: %s", key)
			continue
		}
		mappedScores[questionName] = scores
	}
	// Calculate final sentiment using the default weights.
	finalSentiment := CalculateFinalSentiment(defaultWeights, mappedScores)
	response := SentimentResponse{
		FinalSentiment: finalSentiment,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
		return
	}
}

// -- WebSocket Support --

var upgrader = websocket.Upgrader{
	// Allow any origin for demonstration purposes. Adjust as needed.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// wsHandler upgrades the connection to a WebSocket and sends sentiment updates every 5 minutes.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Immediately send an update, then every 5 minutes
	sendSentimentUpdate(conn)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sendSentimentUpdate(conn)
	}
}

// sendSentimentUpdate calculates the sentiment using defaultMessageScores,
// then sends the result as a JSON message via the WebSocket connection.
func sendSentimentUpdate(conn *websocket.Conn) {
	mappedScores := make(map[string][]float64)
	// Map numeric keys in defaultMessageScores to question strings.
	for key, scores := range defaultMessageScores {
		if questionName, ok := questionMapping[key]; ok {
			mappedScores[questionName] = scores
		}
	}
	finalSentiment := CalculateFinalSentiment(defaultWeights, mappedScores)
	resp := SentimentResponse{
		FinalSentiment: finalSentiment,
	}
	if err := conn.WriteJSON(resp); err != nil {
		log.Println("Error sending WebSocket message:", err)
	}
}

func main() {

	  // 1) Load .env (silently continue if it's missing)
	  if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to env vars")
	  }
	
	// 1) Init DB
	initDB()


	// Load the question mapping from mapping.json.
	data, err := ioutil.ReadFile("mapping.json")
	if err != nil {
		log.Fatalf("Error reading mapping file: %v", err)
	}
	if err := json.Unmarshal(data, &questionMapping); err != nil {
		log.Fatalf("Error parsing mapping file: %v", err)
	}

	// Set up HTTP handlers.
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/sentiment", sentimentHandler)
	// Add the WebSocket endpoint.
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}


/*
{
  _id: ObjectId(),
  source:   "twitter",
  externalId: "1234567890",
  question: "1",
  createdAt: ISODate("2025-04-16T15:00:00Z"),
  fetchedAt: ISODate(),
  author:   "@alice",
  content:  "This coin is amazing because …",
}


{

	metadata: {
    retweetCount: 10,
    likes: 53,
    hashtags: ["#crypto","#DeFi"]
  }
}

*/