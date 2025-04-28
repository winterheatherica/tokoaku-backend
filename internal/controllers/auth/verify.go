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
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
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

	if body.Email == "" || body.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email dan token wajib diisi",
		})
	}

	ctx := context.Background()

	prefix, err := utils.GetVolatileRedisPrefix()
	if err != nil {
		log.Println("Gagal ambil volatile prefix:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal inisialisasi cache",
		})
	}

	redisClient, err := redis.GetRedisClient(prefix)
	if err != nil {
		log.Println("Gagal ambil Redis client:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengakses cache",
		})
	}

	lockKey := "lock:verify:" + body.Email
	ok, err := redisClient.SetNX(ctx, lockKey, "locked", 30*time.Second).Result()
	if err != nil || !ok {
		log.Println("Verifikasi sudah berjalan untuk", body.Email)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"message": "Verifikasi sedang diproses. Coba beberapa saat lagi.",
		})
	}
	defer redisClient.Del(ctx, lockKey)

	redisKey := "verify:" + body.Email
	savedToken, err := redisClient.Get(ctx, redisKey).Result()
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
	passwordPlain, err := redisClient.Get(ctx, passKey).Result()
	if err != nil {
		log.Println("Gagal ambil password plaintext:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal ambil password. Silakan daftar ulang.",
		})
	}

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
		RoleID:       1,
		ProviderID:   1,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		log.Println("Gagal simpan user ke database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan user",
		})
	}

	roleName, err := utils.GetRoleNameByID(int(newUser.RoleID))
	if err != nil {
		log.Println("Gagal ambil nama role:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil role user",
		})
	}

	claims := map[string]interface{}{
		"role": roleName,
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

	_ = redisClient.Del(ctx, redisKey)
	_ = redisClient.Del(ctx, passKey)
	_ = database.DB.Delete(&models.PendingUser{}, "email = ?", body.Email)

	log.Println("Verifikasi berhasil & user dibuat:", userRecord.Email)

	return c.JSON(fiber.Map{
		"message":     "Akun berhasil diverifikasi dan dibuat",
		"customToken": customToken,
		"email":       body.Email,
		"uid":         userRecord.UID,
		"role":        roleName,
	})
}
