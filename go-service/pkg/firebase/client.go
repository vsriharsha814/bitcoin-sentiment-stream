package firebase

import (
    "context"
    "log"
    "os"

    firebase "firebase.google.com/go/v4"
    "cloud.google.com/go/firestore"
    "google.golang.org/api/option"
)

var (
    App    *firebase.App
    client *firestore.Client
)

func Init() {
    ctx := context.Background()
    sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
    app, err := firebase.NewApp(ctx, nil, sa)
    if err != nil {
        log.Fatalf("firebase.NewApp: %v", err)
    }
    App = app

    client, err = app.Firestore(ctx)
    if err != nil {
        log.Fatalf("firestore.NewClient: %v", err)
    }
}

func Client() *firestore.Client {
    return client
}
