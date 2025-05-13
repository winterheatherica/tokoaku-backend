package visitor

import (
	"github.com/gofiber/fiber/v2"

	cloudutil "github.com/winterheatherica/tokoaku-backend/internal/utils/cloudinary"
)

func GetCloudinaryPublicImagePrefix(c *fiber.Ctx) error {
	prefix, err := cloudutil.ResolveCloudinaryPublicPrefix()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"env_key_prefix": prefix,
	})
}
