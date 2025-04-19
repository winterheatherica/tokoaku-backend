package config

import (
	"log"
	"os"
)

var Cloudinary struct {
	CloudName string
	APIKey    string
	APISecret string
}

func LoadCloudinaryConfig() {
	Cloudinary.CloudName = os.Getenv("CLOUDINARY_CLOUD_NAME")
	Cloudinary.APIKey = os.Getenv("CLOUDINARY_API_KEY")
	Cloudinary.APISecret = os.Getenv("CLOUDINARY_API_SECRET")

	log.Println("[CONFIG]: ⚙️  Cloudinary config initialized")
}
