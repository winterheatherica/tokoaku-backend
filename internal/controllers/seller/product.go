package sellercontroller

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/services/redis"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
)

func AddProduct(c *fiber.Ctx) error {
	type ProductInput struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		SellerID    string `json:"seller_id"`
		ProductType string `json:"product_type"`
	}

	var input ProductInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	var productType models.ProductType
	if err := database.DB.Where("name = ?", input.ProductType).First(&productType).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "ProductType not found")
	}

	product := models.Product{
		Name:          input.Name,
		Description:   input.Description,
		SellerID:      input.SellerID,
		ProductTypeID: uint(productType.ID),
	}

	if err := database.DB.Create(&product).Error; err != nil {
		log.Println("[ERROR]: Failed to create product:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create product")
	}

	ctx := context.Background()

	prefix, err := utils.GetVolatileRedisPrefix()
	if err != nil {
		log.Println("[ERROR]: Gagal ambil volatile prefix:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	redisClient, err := redis.GetRedisClient(prefix)
	if err != nil {
		log.Println("[ERROR]: Gagal ambil redis client:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	_ = redisClient.Del(ctx, "product_types")

	return c.JSON(product)
}

func GetAllProductTypes(c *fiber.Ctx) error {
	productTypes, err := utils.GetAllProductTypesFromCacheOrDB()
	if err != nil {
		log.Println("[ERROR]: Failed to fetch product types:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch product types")
	}

	return c.JSON(productTypes)
}
