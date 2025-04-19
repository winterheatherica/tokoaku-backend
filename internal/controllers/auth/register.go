package auth

import (
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
	"github.com/winterheatherica/tokoaku-backend/internal/services/redis"
)

func Register(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body Request
	if err := c.BodyParser(&body); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Data tidak valid",
		})
	}

	log.Println("Cek apakah email sudah terdaftar...")

	if body.Email == "" || body.Password == "" {
		log.Println("Email atau password kosong")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email dan password wajib diisi",
		})
	}

	var existingPending models.PendingUser
	err := database.DB.First(&existingPending, "email = ?", body.Email).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Error saat cek pending_users:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Terjadi kesalahan saat cek email",
		})
	}
	if err == nil {
		redisKey := "verify:" + body.Email
		_, err := redis.Client.Get(redis.Ctx, redisKey).Result()
		if err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Email sudah digunakan dan menunggu verifikasi.",
			})
		}
		_ = database.DB.Delete(&models.PendingUser{}, "email = ?", body.Email)
	}

	cacheKey := "register-check:" + body.Email
	cachedResult, _ := redis.Client.Get(redis.Ctx, cacheKey).Result()
	if cachedResult == "taken" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email sudah terdaftar.",
		})
	}

	var user models.User
	if err := database.DB.First(&user, "email = ?", body.Email).Error; err == nil {
		redis.Client.Set(redis.Ctx, cacheKey, "taken", 10*time.Minute)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email sudah terdaftar.",
		})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		log.Println("Hash error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal hash password",
		})
	}

	log.Println("Password berhasil di-hash")

	passKey := "plainpass:" + body.Email
	err = redis.Client.Set(redis.Ctx, passKey, body.Password, 15*time.Minute).Err()
	if err != nil {
		log.Println("Gagal simpan password plaintext ke Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal simpan password sementara",
		})
	}

	pending := models.PendingUser{
		Email:        body.Email,
		PasswordHash: string(hashed),
	}

	if err := database.DB.Create(&pending).Error; err != nil {
		log.Println("Gagal simpan ke pending_users:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Email sudah digunakan atau gagal simpan",
		})
	}

	log.Println("Data pending user tersimpan")

	token := uuid.NewString()
	redisKey := "verify:" + body.Email

	err = redis.Client.Set(redis.Ctx, redisKey, token, 15*time.Minute).Err()
	if err != nil {
		log.Println("Redis set token error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal simpan token verifikasi",
		})
	}

	log.Printf("Token %s disimpan ke Redis dengan key %s\n", token, redisKey)

	err = email.SendVerificationEmail(body.Email, token)
	if err != nil {
		log.Println("Gagal kirim email verifikasi:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal kirim email verifikasi",
		})
	}

	log.Println("Email verifikasi berhasil dikirim ke", body.Email)
	log.Println("Proses register selesai")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Silakan cek email untuk verifikasi",
		"debug_token": token,
	})
}
