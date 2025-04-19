package firebase

import (
	"context"
	"encoding/json"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var App *firebase.App
var FirebaseAuth *auth.Client

func InitFirebase() {
	raw := os.Getenv("FIREBASE_CREDENTIAL_JSON_RAW")
	if raw == "" {
		log.Fatal("FIREBASE_CREDENTIAL_JSON_RAW kosong! Cek .env kamu")
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &jsonMap); err != nil {
		log.Fatalf("Credential JSON tidak bisa diparse: %v", err)
	}

	credBytes, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatalf("Gagal marshal credential: %v", err)
	}

	opt := option.WithCredentialsJSON(credBytes)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Gagal inisialisasi Firebase: %v", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Gagal inisialisasi Firebase Auth: %v", err)
	}

	App = app
	FirebaseAuth = authClient
	log.Println("[SERVICE]: ⚙️  Firebase Admin SDK connected")
}
