package seller

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func SetVariantCover(c *fiber.Ctx) error {
	variantID := c.Params("id")
	ctx := context.Background()

	var payload struct {
		ImageURL string `json:"image_url"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if payload.ImageURL == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Image URL is required")
	}

	if err := resetCover(variantID); err != nil {
		log.Printf("[DB] ❌ Failed to reset previous cover for variant %s: %v", variantID, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to reset previous cover")
	}

	if err := setCoverImage(variantID, payload.ImageURL); err != nil {
		log.Printf("[DB] ❌ Failed to set new cover image for variant %s: %v", variantID, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to set new cover image")
	}

	var variant models.ProductVariant
	if err := database.DB.Select("product_id").Where("id = ?", variantID).First(&variant).Error; err != nil {
		log.Printf("[DB] ❌ Failed to fetch product_id for variant %s: %v", variantID, err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to refresh cache")
	}

	if err := ClearVariantImageCache(ctx, variantID); err != nil {
		log.Printf("[CACHE] ⚠️ Failed to clear cache for variant %s: %v", variantID, err)
	}

	if err := fetcher.CacheVariantImageFromDB(ctx, variantID); err != nil {
		log.Printf("[CACHE] ⚠️ Failed to refresh variant image cache from DB: %v", err)
	}

	return c.JSON(fiber.Map{
		"message": "Variant cover updated successfully",
	})
}

func resetCover(variantID string) error {
	return database.DB.Model(&models.ProductVariantImage{}).
		Where("product_variant_id = ?", variantID).
		Update("is_variant_cover", false).Error
}

func setCoverImage(variantID string, imageURL string) error {
	return database.DB.Model(&models.ProductVariantImage{}).
		Where("product_variant_id = ? AND image_url = ?", variantID, imageURL).
		Update("is_variant_cover", true).Error
}

func ClearVariantImageCache(ctx context.Context, variantID string) error {
	rdb, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("product_variant_images:%s", variantID)
	return rdb.Del(ctx, key).Err()
}
