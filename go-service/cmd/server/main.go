package main

import (
    "io/ioutil"
    "log"
    "net/http"
    "os"
	"strings"

    "github.com/joho/godotenv"

    "github.com/cosmic-hash/CryptoPulse/pkg/config"
    "github.com/cosmic-hash/CryptoPulse/pkg/db"
    handlers "github.com/cosmic-hash/CryptoPulse/pkg/handler"
	"github.com/cosmic-hash/CryptoPulse/pkg/firebase"
	openai "github.com/cosmic-hash/CryptoPulse/pkg/openai"
	 
)

func main() {
    // 1) Load .env if it exists
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, falling back to env vars")
    }
	// DEBUG: verify the key is there
    //log.Printf("→ OPENAI_API_KEY=%q", os.Getenv("OPENAI_API_KEY"))

    // 1.a) Initialize the OpenAI client now that the .env is loaded
    openai.InitClient()
    // 1.a) Debug: print out the DATABASE_URL you're using
    // log.Printf("→ DATABASE_URL=%q", os.Getenv("DATABASE_URL"))

    // 2) Init DB (will log fatal if it still can’t connect)
    db.InitDB()
	firebase.Init()

    // 3) Load question mapping
    mappingPath := os.Getenv("QUESTION_MAPPING_FILE")
    if mappingPath == "" {
        mappingPath = "mapping.json"
    }
    data, err := ioutil.ReadFile(mappingPath)
    if err != nil {
        log.Fatalf("Error reading mapping file %q: %v", mappingPath, err)
    }
    if err := config.LoadQuestionMapping(data); err != nil {
        log.Fatalf("Error parsing mapping file: %v", err)
    }

// // /alerts → both list (GET) and create (POST)
// http.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		handlers.ListAlertsHandler(w, r)
// 	case http.MethodPost:
// 		handlers.CreateAlertHandler(w, r)
// 	default:
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 	}
// })

// /alerts/{id} → delete only
http.HandleFunc("/alerts/", func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// the path is "/alerts/{id}"
	// so everything after "/alerts/" is the id
	id := strings.TrimPrefix(r.URL.Path, "/alerts/")
	// inject into the Request as a query param for simplicity
	q := r.URL.Query()
	q.Set("id", id)
	r.URL.RawQuery = q.Encode()

	handlers.DeleteAlertHandler(w, r)
})

http.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        handlers.ListAlertsHandler(w, r)

    case http.MethodPost:
        handlers.CreateAlertHandler(w, r)

    case http.MethodDelete:
        // we’ll delete by coinId query param, not Firestore ID
        coinStr := r.URL.Query().Get("coinId")
        if coinStr == "" {
            http.Error(w, "coinId query required", http.StatusBadRequest)
            return
        }
        // inject coinId into the request context so your handler can read it:
        q := r.URL.Query()
        q.Set("coinId", coinStr)
        r.URL.RawQuery = q.Encode()

        handlers.DeleteAlertHandler(w, r)

    default:
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
})

    // 4) Register HTTP & WebSocket handlers
    http.HandleFunc("/", handlers.HelloHandler)
    http.HandleFunc("/sentiment", handlers.SentimentHandler)
    http.HandleFunc("/ws", handlers.WSHandler)
	http.HandleFunc("/aggregate", handlers.AggregateHandler)
	http.HandleFunc("/explain", handlers.ExplainSentimentHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🟢 Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
