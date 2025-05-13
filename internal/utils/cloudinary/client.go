package cloudinary

import (
	cldlib "github.com/cloudinary/cloudinary-go/v2"
)

func GetCloudinaryClient(prefix string) (*cldlib.Cloudinary, error) {
	cfg, err := GetCloudinaryConfig(prefix)
	if err != nil {
		return nil, err
	}
	return cldlib.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
}

func GetCloudinaryClientFromConfig(cfg *CloudinaryConfig) (*cldlib.Cloudinary, error) {
	return cldlib.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
}
