package visitor

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetProductTypes(c *fiber.Ctx) error {
	ctx := context.Background()

	types, err := fetcher.GetAllProductTypes(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data product types")
	}

	return c.JSON(types)
}
