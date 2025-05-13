package utils

// // OpenUploadedFile opens the file from fiber form header
// func OpenUploadedFile(fileHeader *multipart.FileHeader) (multipart.File, error) {
// 	return fileHeader.Open()
// }

// // IsImageSquare checks if image has 1:1 aspect ratio
// func IsImageSquare(img image.Image) bool {
// 	bounds := img.Bounds()
// 	return bounds.Dx() == bounds.Dy()
// }

// // ResizeIfNeeded resizes the image to 512x512 if larger
// func ResizeIfNeeded(img image.Image) image.Image {
// 	if img.Bounds().Dx() > 512 {
// 		return resize.Resize(512, 512, img, resize.Lanczos3)
// 	}
// 	return img
// }

// // IsSupportedImageFormat checks if the format is jpeg/jpg or png
// func IsSupportedImageFormat(format string) bool {
// 	switch format {
// 	case "jpeg", "jpg", "png":
// 		return true
// 	default:
// 		return false
// 	}
// }

// // EncodeImage encodes image.Image to *bytes.Buffer based on format
// func EncodeImage(img image.Image, format string) (*bytes.Buffer, error) {
// 	var buf bytes.Buffer
// 	var err error

// 	switch format {
// 	case "jpeg", "jpg":
// 		err = jpeg.Encode(&buf, img, nil)
// 	case "png":
// 		err = png.Encode(&buf, img)
// 	default:
// 		err = fmt.Errorf("unsupported image format")
// 	}

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &buf, nil
// }
