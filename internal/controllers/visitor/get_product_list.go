package visitor

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetVisitorProductList(c *fiber.Ctx) error {
	response, err := fetcher.FetchProductCardList()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil produk")
	}
	return c.JSON(response)
}
