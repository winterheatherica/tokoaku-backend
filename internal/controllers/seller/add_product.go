package seller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
	cloudutil "github.com/winterheatherica/tokoaku-backend/internal/utils/cloudinary"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
)

type AddProductRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	ProductType   string `json:"product_type"`
	ProductForm   string `json:"product_form"`
	ImageCoverURL string `json:"image_cover_url"`
}

func AddProduct(c *fiber.Ctx) error {
	var req AddProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input",
		})
	}

	sellerID := c.Locals("uid").(string)

	productTypes, err := fetcher.GetAllProductTypes(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load product types")
	}

	var pt *models.ProductType
	for _, t := range productTypes {
		if t.Name == req.ProductType {
			pt = &t
			break
		}
	}
	if pt == nil {
		return fiber.NewError(fiber.StatusNotFound, "Product type not found")
	}

	productForms, err := fetcher.GetAllProductForms(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load product forms")
	}

	var pf *models.ProductForm
	for _, form := range productForms {
		if form.Form == req.ProductForm {
			pf = &form
			break
		}
	}
	if pf == nil {
		return fiber.NewError(fiber.StatusNotFound, "Product form not found")
	}

	prefix, err := cloudutil.ResolveCloudinaryPublicPrefix()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cloudinary prefix gagal diambil: "+err.Error())
	}

	var cs models.CloudService
	if err := database.DB.Where("env_key_prefix = ?", prefix).First(&cs).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cloudinary service tidak ditemukan di DB")
	}

	productID := uuid.New().String()
	maxAttempts := 5

	for i := 0; i < maxAttempts; i++ {
		slug, err := utils.GenerateUniqueSlug(req.Name)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Gagal membuat slug")
		}

		newProduct := models.Product{
			ID:            productID,
			Name:          req.Name,
			Description:   req.Description,
			SellerID:      sellerID,
			ProductTypeID: pt.ID,
			CloudImageID:  cs.ID,
			ProductFormID: pf.ID,
			ImageCoverURL: req.ImageCoverURL,
			Slug:          slug,
		}

		if err := database.DB.Create(&newProduct).Error; err != nil {
			if database.DB.Where("slug = ?", slug).First(&models.Product{}).Error == nil {
				continue
			}
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create product")
		}

		return c.Status(fiber.StatusCreated).JSON(newProduct)
	}

	return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate unique slug after multiple attempts")
}
