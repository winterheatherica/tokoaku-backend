package customer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

func LetReview(c *fiber.Ctx) error {
	userUID := c.Locals("uid").(string)
	productSlug := c.Params("product_slug")

	var product models.Product
	if err := database.DB.Where("slug = ?", productSlug).First(&product).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Produk tidak ditemukan")
	}

	var allVariantIDs []string
	if err := database.DB.Model(&models.ProductVariant{}).
		Where("product_id = ?", product.ID).
		Pluck("id", &allVariantIDs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil varian produk")
	}

	var orders []models.Order
	if err := database.DB.Where("customer_id = ?", userUID).Find(&orders).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil order")
	}

	var paidOrderIDs []uint
	for _, o := range orders {
		var logs []models.OrderLog
		if err := database.DB.
			Where("order_id = ?", o.ID).
			Order("created_at DESC").
			Find(&logs).Error; err == nil && len(logs) > 0 && logs[0].StatusID == 12 {
			paidOrderIDs = append(paidOrderIDs, o.ID)
		}
	}

	if len(paidOrderIDs) == 0 {
		return c.JSON(fiber.Map{
			"can_review":        false,
			"message":           "Kamu belum pernah melakukan pembelian produk ini",
			"eligible_variants": []fiber.Map{},
		})
	}

	var shippingIDs []uint
	if err := database.DB.Model(&models.OrderShipping{}).
		Where("order_id IN ?", paidOrderIDs).
		Pluck("id", &shippingIDs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil shipping")
	}

	var boughtVariantIDs []string
	if err := database.DB.Model(&models.OrderItem{}).
		Where("order_shipping_id IN ?", shippingIDs).
		Pluck("product_variant_id", &boughtVariantIDs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil item")
	}

	var eligibleVariants []models.ProductVariant
	if err := database.DB.
		Where("id IN ?", boughtVariantIDs).
		Where("product_id = ?", product.ID).
		// Where("id NOT IN (?)", database.DB.Table("reviews").Select("product_variant_id").Where("customer_id = ?", userUID)).
		Find(&eligibleVariants).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memfilter varian yang belum direview")
	}

	response := []fiber.Map{}
	for _, v := range eligibleVariants {
		response = append(response, fiber.Map{
			"id":           v.ID,
			"variant_name": v.VariantName,
			"slug":         v.Slug,
		})
	}

	if len(response) == 0 {
		return c.JSON(fiber.Map{
			"can_review":        false,
			"message":           "Kamu sudah memberikan review untuk semua varian yang pernah kamu beli.",
			"eligible_variants": []fiber.Map{},
		})
	}

	return c.JSON(fiber.Map{
		"can_review":        true,
		"message":           "Pilih varian yang ingin kamu review",
		"eligible_variants": response,
	})
}
