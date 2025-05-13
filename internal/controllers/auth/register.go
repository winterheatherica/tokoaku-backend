package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/services/email"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func Register(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body Request
	if err := c.BodyParser(&body); err != nil || body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email dan password wajib diisi",
		})
	}

	ctx := context.Background()

	redisClient, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		log.Println("Redis init error:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal inisialisasi cache")
	}

	var existingPending models.PendingUser
	err = database.DB.First(&existingPending, "email = ?", body.Email).Error
	if err == nil {
		if _, err := redisClient.Get(ctx, "verify:"+body.Email).Result(); err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Email sudah digunakan dan menunggu verifikasi.",
			})
		}
		_ = database.DB.Delete(&models.PendingUser{}, "email = ?", body.Email)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Error saat cek pending_users:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Terjadi kesalahan saat cek email")
	}

	cacheKey := "register-check:" + body.Email
	if val, _ := redisClient.Get(ctx, cacheKey).Result(); val == "taken" {
		return fiber.NewError(fiber.StatusBadRequest, "Email sudah terdaftar.")
	}

	var user models.User
	if err := database.DB.First(&user, "email = ?", body.Email).Error; err == nil {
		_ = redisClient.Set(ctx, cacheKey, "taken", 10*time.Minute).Err()
		return fiber.NewError(fiber.StatusBadRequest, "Email sudah terdaftar.")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		log.Println("Hash error:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal hash password")
	}

	if err := redisClient.Set(ctx, "plainpass:"+body.Email, body.Password, 15*time.Minute).Err(); err != nil {
		log.Println("Gagal simpan plaintext password:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal simpan password sementara")
	}

	pending := models.PendingUser{
		Email:        body.Email,
		PasswordHash: string(hashed),
	}
	if err := database.DB.Create(&pending).Error; err != nil {
		log.Println("Gagal simpan pending user:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Email sudah digunakan atau gagal simpan")
	}

	token := uuid.NewString()
	if err := redisClient.Set(ctx, "verify:"+body.Email, token, 15*time.Minute).Err(); err != nil {
		log.Println("Gagal simpan token verifikasi:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal simpan token verifikasi")
	}

	if err := email.SendVerificationEmail(body.Email, token); err != nil {
		log.Println("Gagal kirim email verifikasi:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal kirim email verifikasi")
	}

	log.Printf("✅ Register berhasil untuk %s, token: %s\n", body.Email, token)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Silakan cek email untuk verifikasi",
		"debug_token": token,
	})
}
