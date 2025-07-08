package customer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type AddAddressInput struct {
	AddressLine string  `json:"address_line"`
	City        string  `json:"city"`
	Province    string  `json:"province"`
	PostalCode  string  `json:"postal_code"`
	IsActive    bool    `json:"is_active"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

func AddAddress(c *fiber.Ctx) error {
	customerID := c.Locals("uid").(string)

	var input AddAddressInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if input.AddressLine == "" || input.City == "" || input.Province == "" || input.PostalCode == "" {
		return fiber.NewError(fiber.StatusBadRequest, "all fields are required")
	}

	if input.Latitude == 0 || input.Longitude == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "location (latitude & longitude) is required")
	}

	// Nonaktifkan semua alamat aktif sebelumnya jika alamat ini dijadikan utama
	if input.IsActive {
		if err := database.DB.Model(&models.Address{}).
			Where("user_id = ? AND is_active = true", customerID).
			Update("is_active", false).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to update address status")
		}
	}

	address := &models.Address{
		ID:          uuid.New(),
		UserID:      customerID,
		AddressLine: input.AddressLine,
		City:        input.City,
		Province:    input.Province,
		PostalCode:  input.PostalCode,
		Latitude:    input.Latitude,
		Longitude:   input.Longitude,
		IsActive:    input.IsActive,
	}

	if err := database.DB.Create(address).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to save address")
	}

	return c.JSON(fiber.Map{
		"message": "address added successfully",
		"address": address,
	})
}
