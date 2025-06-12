package visitor

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

func GetvisitorProductReferenceData(c *fiber.Ctx) error {
	ctx := context.Background()

	forms, err := fetcher.GetAllProductForms(ctx)
	if err != nil {
		log.Println("❌ Gagal ambil product forms:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil daftar product forms",
		})
	}

	types, err := fetcher.GetAllProductTypes(ctx)
	if err != nil {
		log.Println("❌ Gagal ambil product types:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil daftar product types",
		})
	}

	labels, err := fetcher.GetAllCategoryLabels(ctx)
	if err != nil {
		log.Println("❌ Gagal ambil category labels:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil daftar category labels",
		})
	}

	categories, err := fetcher.GetAllCategories(ctx)
	if err != nil {
		log.Println("❌ Gagal ambil categories:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil daftar categories",
		})
	}

	log.Printf("✅ Jumlah labels ditemukan: %d", len(labels))
	log.Printf("✅ Jumlah categories ditemukan: %d", len(categories))

	categoryMap := make(map[uint][]fiber.Map)
	for _, cat := range categories {
		categoryMap[cat.CategoryLabelID] = append(categoryMap[cat.CategoryLabelID], fiber.Map{
			"id":   cat.ID,
			"name": cat.Name,
			"slug": cat.Slug,
			"code": cat.Code,
		})
	}

	for labelID, cats := range categoryMap {
		log.Printf("📦 Label ID %d punya %d kategori", labelID, len(cats))
	}

	var groupedLabels []fiber.Map
	for _, label := range labels {
		cats := categoryMap[label.ID]
		if cats == nil {
			cats = []fiber.Map{}
			log.Printf("⚠️ Label ID %d (%s) tidak punya kategori", label.ID, label.Name)
		} else {
			log.Printf("✅ Label ID %d (%s) akan dikirim dengan %d kategori", label.ID, label.Name, len(cats))
		}

		groupedLabels = append(groupedLabels, fiber.Map{
			"id":         label.ID,
			"name":       label.Name,
			"slug":       label.Slug,
			"categories": cats,
		})

	}

	return c.JSON(fiber.Map{
		"product_forms":   forms,
		"product_types":   types,
		"category_labels": groupedLabels,
	})
}
