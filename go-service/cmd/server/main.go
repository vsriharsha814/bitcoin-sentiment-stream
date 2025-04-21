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
	"github.com/cosmic-hash/CryptoPulse/pkg/handler"
	
)

func main() {
	// 1) Load .env (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to env vars")
	}

	// 2) Init DB
	db.InitDB()

	// 3) Load the question mapping from mapping.json
	path := "mapping.json"
	if envPath := os.Getenv("QUESTION_MAPPING_FILE"); envPath != "" {
		path = envPath
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading mapping file: %v", err)
	}
	if err := config.LoadQuestionMapping(data); err != nil {
		log.Fatalf("Error parsing mapping file: %v", err)
	}

	// 4) Register handlers
	http.HandleFunc("/", handlers.HelloHandler)
	http.HandleFunc("/sentiment", handlers.SentimentHandler)
	http.HandleFunc("/ws", handlers.WSHandler)

	fmt.Println("Server is listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
