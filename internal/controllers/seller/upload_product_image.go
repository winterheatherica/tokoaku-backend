package seller

import (
	"github.com/gofiber/fiber/v2"
	cloudutil "github.com/winterheatherica/tokoaku-backend/internal/utils/cloudinary"
)

func UploadProductImage(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "File not found")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open uploaded file")
	}
	defer file.Close()

	buf, err := cloudutil.ProcessSquareImage(file)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	prefix, err := cloudutil.ResolveCloudinaryPublicPrefix()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	url, err := cloudutil.UploadBufferToCloudinary(prefix, "products", buf)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"secure_url": url,
	})
}
