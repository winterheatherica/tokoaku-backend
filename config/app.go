package config

import (
	"log"
	"os"
)

var App struct {
	Port            string
	Env             string
	FrontendBaseURL string
	BackendBaseURL  string
	MLBaseURL       string
}

func LoadAppConfig() {
	App.Port = os.Getenv("PORT")
	App.Env = os.Getenv("APP_ENV")
	App.FrontendBaseURL = os.Getenv("FRONTEND_BASE_URL")
	App.BackendBaseURL = os.Getenv("BACKEND_BASE_URL")
	App.MLBaseURL = os.Getenv("MACHINE_LEARNING_BASE_URL")

	log.Println("[CONFIG]: ⚙️  App config initialized")
}
