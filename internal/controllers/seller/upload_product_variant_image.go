package seller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	cloudutil "github.com/winterheatherica/tokoaku-backend/internal/utils/cloudinary"
)

func UploadProductVariantImage(c *fiber.Ctx) error {
	productVariantID := c.Params("id")

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

	url, err := cloudutil.UploadBufferToCloudinary(prefix, "product-variants", buf)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	imageRecord := models.ProductVariantImage{
		ProductVariantID: productVariantID,
		ImageURL:         url,
		CloudImageID:     MustGetCloudImageID(prefix),
		IsVariantCover:   false,
	}
	if err := database.DB.Create(&imageRecord).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save image record")
	}

	return c.JSON(fiber.Map{
		"message": "Variant image uploaded successfully",
		"image":   imageRecord,
	})
}

func MustGetCloudImageID(prefix string) uint {
	var cs models.CloudService
	_ = database.DB.Where("env_key_prefix = ?", prefix).First(&cs).Error
	return cs.ID
}
