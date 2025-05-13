package visitor

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetProductForms(c *fiber.Ctx) error {
	ctx := context.Background()

	forms, err := fetcher.GetAllProductForms(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data product forms")
	}

	return c.JSON(forms)
}
