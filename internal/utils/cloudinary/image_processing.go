package cloudinary

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"

	"github.com/nfnt/resize"
)

func ProcessSquareImage(file multipart.File) (*bytes.Buffer, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("invalid image file: %w", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width != height {
		return nil, fmt.Errorf("image must be 1:1 aspect ratio")
	}

	if width > 512 {
		img = resize.Resize(512, 512, img, resize.Lanczos3)
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, nil)
	case "png":
		err = png.Encode(&buf, img)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return &buf, nil
}

func ProcessBannerImage(file multipart.File) (*bytes.Buffer, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("invalid image file: %w", err)
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	targetRatio := 2100.0 / 900.0
	currentRatio := float64(width) / float64(height)

	var cropWidth, cropHeight int
	if currentRatio > targetRatio {
		cropHeight = height
		cropWidth = int(float64(height) * targetRatio)
	} else {
		cropWidth = width
		cropHeight = int(float64(width) / targetRatio)
	}
	offsetX := (width - cropWidth) / 2
	offsetY := (height - cropHeight) / 2

	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(offsetX, offsetY, offsetX+cropWidth, offsetY+cropHeight))

	if cropWidth > 2100 {
		croppedImg = resize.Resize(2100, 0, croppedImg, resize.Lanczos3)
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, croppedImg, nil)
	case "png":
		err = png.Encode(&buf, croppedImg)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return &buf, nil
}
