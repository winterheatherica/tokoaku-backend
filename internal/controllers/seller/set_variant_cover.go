package seller

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/persistent"
)

func SetVariantCover(c *fiber.Ctx) error {
	variantID := c.Params("id")

	var payload struct {
		ImageURL string `json:"image_url"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if payload.ImageURL == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Image URL is required")
	}

	if err := ResetCover(variantID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to reset previous cover")
	}

	if err := SetCoverImage(variantID, payload.ImageURL); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to set new cover image")
	}

	ctx := context.Background()
	ClearVariantCoverCache(ctx, variantID)
	if _, err := fetcher.GetVariantCoverImage(ctx, variantID); err != nil {
		log.Printf("⚠️ Gagal refresh cache cover variant %s: %v", variantID, err)
	}

	return c.JSON(fiber.Map{
		"message": "Variant cover updated successfully",
	})
}

func ResetCover(variantID string) error {
	return database.DB.
		Model(&models.ProductVariantImage{}).
		Where("product_variant_id = ?", variantID).
		Update("is_variant_cover", false).Error
}

func SetCoverImage(variantID string, imageURL string) error {
	return database.DB.
		Model(&models.ProductVariantImage{}).
		Where("product_variant_id = ? AND image_url = ?", variantID, imageURL).
		Update("is_variant_cover", true).Error
}

func ClearVariantCoverCache(ctx context.Context, variantID string) {
	rdb, err := persistent.GetPersistentRedisClient(ctx)
	if err == nil {
		_ = rdb.Del(ctx, CacheKeyVariantCover(variantID)).Err()
	}
}

func CacheKeyVariantCover(variantID string) string {
	return fmt.Sprintf("variant:cover:%s", variantID)
}
