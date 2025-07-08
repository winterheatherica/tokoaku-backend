package seller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type AddBankAccountRequest struct {
	BankID        uint   `json:"bank_id" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required"`
	AccountName   string `json:"account_name" validate:"required"`
	IsActive      bool   `json:"is_active"`
}

func AddBankAccount(c *fiber.Ctx) error {
	var body AddBankAccountRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Data tidak valid")
	}

	uid := c.Locals("uid").(string)

	bankAccount := models.BankAccount{
		ID:            uuid.New(),
		UserID:        uid,
		BankID:        body.BankID,
		AccountNumber: body.AccountNumber,
		AccountName:   body.AccountName,
		IsActive:      body.IsActive,
		CreatedAt:     time.Now(),
	}

	if err := database.DB.Create(&bankAccount).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan akun bank")
	}

	return c.JSON(fiber.Map{
		"message": "Akun bank berhasil ditambahkan",
		"data":    bankAccount,
	})
}

func SetActiveBankAccount(c *fiber.Ctx) error {
	userID := c.Locals("uid").(string)
	accountID := c.Params("id")

	// Pastikan account milik user yang sedang login
	var target models.BankAccount
	if err := database.DB.
		Where("id = ? AND user_id = ?", accountID, userID).
		First(&target).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Akun bank tidak ditemukan atau bukan milik Anda")
	}

	// Nonaktifkan semua akun lain milik user
	if err := database.DB.
		Model(&models.BankAccount{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menonaktifkan akun-akun lain")
	}

	// Aktifkan akun yang dipilih
	if err := database.DB.
		Model(&models.BankAccount{}).
		Where("id = ?", accountID).
		Update("is_active", true).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengaktifkan akun bank")
	}

	return c.JSON(fiber.Map{
		"message": "Akun bank berhasil diatur sebagai aktif",
	})
}
