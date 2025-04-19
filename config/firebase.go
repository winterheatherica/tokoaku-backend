package config

import (
	"log"
	"os"
)

var Firebase struct {
	ProjectID     string
	ClientEmail   string
	PrivateKey    string
	Credential    string
	CredentialRaw string
}

func LoadFirebaseConfig() {
	Firebase.ProjectID = os.Getenv("FIREBASE_PROJECT_ID")
	Firebase.ClientEmail = os.Getenv("FIREBASE_CLIENT_EMAIL")
	Firebase.PrivateKey = os.Getenv("FIREBASE_PRIVATE_KEY")
	Firebase.Credential = os.Getenv("FIREBASE_CREDENTIAL_JSON")
	Firebase.CredentialRaw = os.Getenv("FIREBASE_CREDENTIAL_JSON_RAW")

	log.Println("[CONFIG]: ⚙️  Firebase config initialized")
}
