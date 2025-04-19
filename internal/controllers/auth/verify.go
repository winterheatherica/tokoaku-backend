package auth

import (
	"context"
	"log"
	"time"

	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
	"github.com/winterheatherica/tokoaku-backend/internal/services/firebase"
	"github.com/winterheatherica/tokoaku-backend/internal/services/redis"
)

func VerifyToken(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}

	var body Request
	if err := c.BodyParser(&body); err != nil {
		log.Println("BodyParser error (verify):", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Data tidak valid",
		})
	}

	lockKey := "lock:verify:" + body.Email
	ok, err := redis.Client.SetNX(redis.Ctx, lockKey, "locked", 30*time.Second).Result()
	if err != nil || !ok {
		log.Println("Verifikasi sudah berjalan untuk", body.Email)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"message": "Verifikasi sedang diproses. Coba beberapa saat lagi.",
		})
	}
	defer redis.Client.Del(redis.Ctx, lockKey)

	log.Printf("Verifikasi untuk %s dengan token %s\n", body.Email, body.Token)

	if body.Email == "" || body.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email dan token wajib diisi",
		})
	}

	redisKey := "verify:" + body.Email
	savedToken, err := redis.Client.Get(redis.Ctx, redisKey).Result()
	if err != nil || savedToken != body.Token {
		log.Println("Token verifikasi salah atau kadaluarsa")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token tidak valid atau kadaluarsa",
		})
	}

	var pending models.PendingUser
	if err := database.DB.First(&pending, "email = ?", body.Email).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Email belum terdaftar",
		})
	}

	passKey := "plainpass:" + body.Email
	passwordPlain, err := redis.Client.Get(redis.Ctx, passKey).Result()
	if err != nil {
		log.Println("Gagal ambil password plaintext:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal ambil password. Silakan daftar ulang.",
		})
	}

	ctx := context.Background()
	authClient, err := firebase.App.Auth(ctx)
	if err != nil {
		log.Println("Firebase Auth init error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal inisialisasi Firebase",
		})
	}

	userRecord, err := authClient.CreateUser(ctx, (&firebaseauth.UserToCreate{}).
		Email(body.Email).
		Password(passwordPlain))
	if err != nil {
		log.Println("Gagal buat akun Firebase:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat akun Firebase",
		})
	}

	newUser := models.User{
		ID:           userRecord.UID,
		Email:        pending.Email,
		PasswordHash: &pending.PasswordHash,
		Role:         0,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		log.Println("Gagal simpan user ke database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan user",
		})
	}

	claims := map[string]interface{}{
		"role": newUser.Role,
	}
	err = firebase.FirebaseAuth.SetCustomUserClaims(ctx, newUser.ID, claims)
	if err != nil {
		log.Println("Gagal set custom claims:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyetel role user",
		})
	}

	time.Sleep(500 * time.Millisecond)

	customToken, err := firebase.FirebaseAuth.CustomToken(ctx, newUser.ID)
	if err != nil {
		log.Println("Gagal buat custom token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat token login",
		})
	}

	_ = redis.Client.Del(redis.Ctx, redisKey)
	_ = redis.Client.Del(redis.Ctx, passKey)
	_ = database.DB.Delete(&models.PendingUser{}, "email = ?", body.Email)

	log.Println("Verifikasi berhasil dan user dibuat:", userRecord.Email)

	return c.JSON(fiber.Map{
		"message":     "Akun berhasil diverifikasi dan dibuat",
		"customToken": customToken,
		"email":       body.Email,
		"uid":         userRecord.UID,
		"role":        newUser.Role,
	})
}
