package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"

    "github.com/cosmic-hash/CryptoPulse/pkg/config"
    "github.com/cosmic-hash/CryptoPulse/pkg/db"
    handlers "github.com/cosmic-hash/CryptoPulse/pkg/handler"
)

func main() {
    // 1) Load .env if it exists
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, falling back to env vars")
    }

    // 1.a) Debug: print out the DATABASE_URL you're using
    // log.Printf("→ DATABASE_URL=%q", os.Getenv("DATABASE_URL"))

    // 2) Init DB (will log fatal if it still can’t connect)
    db.InitDB()

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

    // 4) Register HTTP & WebSocket handlers
    http.HandleFunc("/", handlers.HelloHandler)
    http.HandleFunc("/sentiment", handlers.SentimentHandler)
    http.HandleFunc("/ws", handlers.WSHandler)

    fmt.Println("Server is listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
