// pkg/openai/client.go
package openai

import (
    "log"
    "os"

    "github.com/openai/openai-go"
    "github.com/openai/openai-go/option"
)

// ChatClient is the shared OpenAI client for your service.
var ChatClient openai.Client

// InitClient must be called after .env is loaded, so OPENAI_API_KEY is present.
func InitClient() {
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENAI_API_KEY is not set")
    }
    // openai.NewClient returns an openai.Client (not *openai.Client)
    ChatClient = openai.NewClient(
        option.WithAPIKey(apiKey),
    )
}