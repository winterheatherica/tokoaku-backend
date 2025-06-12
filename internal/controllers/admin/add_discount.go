package admin

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
	cloudutil "github.com/winterheatherica/tokoaku-backend/internal/utils/cloudinary"
)

func AddDiscount(c *fiber.Ctx) error {
	db := database.DB

	name := c.FormValue("name")

	valueTypeIDStr := c.FormValue("value_type_id")
	valueTypeID, err := strconv.Atoi(valueTypeIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid value_type_id")
	}

	valueStr := c.FormValue("value")
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid value")
	}

	sponsorIDStr := c.FormValue("sponsor_id")
	sponsorID, err := strconv.Atoi(sponsorIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid sponsor_id")
	}

	startAtStr := c.FormValue("start_at")
	startAt, err := time.Parse(time.RFC3339, startAtStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid start_at format (use RFC3339)")
	}

	endAtStr := c.FormValue("end_at")
	endAt, err := time.Parse(time.RFC3339, endAtStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid end_at format (use RFC3339)")
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Image not found")
	}
	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open uploaded image")
	}
	defer file.Close()

	buf, err := cloudutil.ProcessBannerImage(file)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	prefix, err := cloudutil.ResolveCloudinaryPublicPrefix()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	url, err := cloudutil.UploadBufferToCloudinary(prefix, "discounts", buf)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	slugText := utils.SlugifyText(name)

	discount := models.Discount{
		Name:          name,
		ValueTypeID:   uint(valueTypeID),
		Value:         uint(value),
		SponsorID:     uint(sponsorID),
		StartAt:       startAt,
		EndAt:         endAt,
		Slug:          slugText,
		ImageCoverURL: url,
		CloudImageID:  0,
	}

	if err := db.Create(&discount).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create discount: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Discount created successfully",
		"data":    discount,
	})
}
