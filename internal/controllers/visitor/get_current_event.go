package visitor

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetDiscountLimit(c *fiber.Ctx) error {
	ctx := context.Background()

	currentEvent, err := fetcher.GetCurrentEvent(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Tidak ada event aktif saat ini")
	}

	return c.JSON(fiber.Map{
		"discount_limit": currentEvent.EventType.DiscountLimit,
	})
}
