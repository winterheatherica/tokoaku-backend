package config

import (
	"log"
	"os"
)

var Email struct {
	APIKey string
	Sender string
}

func LoadResendConfig() {
	Email.APIKey = os.Getenv("RESEND_API_KEY")
	Email.Sender = os.Getenv("RESEND_SENDER_EMAIL")

	log.Println("[CONFIG]: ⚙️  Resend config initialized")
}
