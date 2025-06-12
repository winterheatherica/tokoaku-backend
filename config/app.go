package config

import (
	"log"
	"os"
)

var App struct {
	Port                   string
	PlatformName           string
	Env                    string
	FrontendBaseURL        string
	BackendBaseURL         string
	MachineLearningBaseURL string
}

func LoadAppConfig() {
	App.Port = os.Getenv("PORT")
	App.PlatformName = os.Getenv("PLATFORM_NAME")
	App.Env = os.Getenv("APP_ENV")
	App.FrontendBaseURL = os.Getenv("FRONTEND_BASE_URL")
	App.BackendBaseURL = os.Getenv("BACKEND_BASE_URL")
	App.MachineLearningBaseURL = os.Getenv("MACHINE_LEARNING_BASE_URL")

	log.Println("[CONFIG]: ⚙️  App config initialized")
}
