package cloudinary

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadBufferToCloudinary(prefix string, folder string, buf *bytes.Buffer) (string, error) {
	cld, err := GetCloudinaryClient(prefix)
	if err != nil {
		return "", fmt.Errorf("cloudinary client error: %w", err)
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), buf, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	return uploadResult.SecureURL, nil
}
