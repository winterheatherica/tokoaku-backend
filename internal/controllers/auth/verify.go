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
	"github.com/winterheatherica/tokoaku-backend/internal/utils/fetcher"
	"github.com/winterheatherica/tokoaku-backend/internal/utils/redis/volatile"
)

func VerifyToken(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}

	var body Request
	if err := c.BodyParser(&body); err != nil || body.Email == "" || body.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email dan token wajib diisi",
		})
	}

	ctx := context.Background()

	redisClient, err := volatile.GetVolatileRedisClient(ctx)
	if err != nil {
		log.Println("Redis init error:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal inisialisasi cache")
	}

	lockKey := "lock:verify:" + body.Email
	if ok, _ := redisClient.SetNX(ctx, lockKey, "locked", 30*time.Second).Result(); !ok {
		return fiber.NewError(fiber.StatusTooManyRequests, "Verifikasi sedang diproses. Coba lagi.")
	}
	defer redisClient.Del(ctx, lockKey)

	savedToken, err := redisClient.Get(ctx, "verify:"+body.Email).Result()
	if err != nil || savedToken != body.Token {
		return fiber.NewError(fiber.StatusUnauthorized, "Token tidak valid atau kadaluarsa")
	}

	var pending models.PendingUser
	if err := database.DB.First(&pending, "email = ?", body.Email).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Email belum terdaftar")
	}

	passwordPlain, err := redisClient.Get(ctx, "plainpass:"+body.Email).Result()
	if err != nil {
		log.Println("Gagal ambil password plaintext:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal ambil password. Silakan daftar ulang.")
	}

	authClient, err := firebase.App.Auth(ctx)
	if err != nil {
		log.Println("Firebase Auth init error:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal inisialisasi Firebase")
	}

	userRecord, err := authClient.CreateUser(ctx, (&firebaseauth.UserToCreate{}).
		Email(body.Email).
		Password(passwordPlain))
	if err != nil {
		log.Println("Gagal buat akun Firebase:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal membuat akun Firebase")
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
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan user")
	}

	roleName, err := fetcher.GetRoleNameByID(int(newUser.RoleID))
	if err != nil {
		log.Println("Gagal ambil nama role:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil role user")
	}

	claims := map[string]interface{}{"role": roleName}
	if err := firebase.FirebaseAuth.SetCustomUserClaims(ctx, newUser.ID, claims); err != nil {
		log.Println("Gagal set custom claims:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyetel role user")
	}

	time.Sleep(500 * time.Millisecond)

	customToken, err := firebase.FirebaseAuth.CustomToken(ctx, newUser.ID)
	if err != nil {
		log.Println("Gagal buat custom token:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal membuat token login")
	}

	_ = redisClient.Del(ctx, "verify:"+body.Email)
	_ = redisClient.Del(ctx, "plainpass:"+body.Email)
	_ = database.DB.Delete(&models.PendingUser{}, "email = ?", body.Email)

	log.Printf("✅ Verifikasi berhasil & user dibuat: %s\n", body.Email)

	return c.JSON(fiber.Map{
		"message":     "Akun berhasil diverifikasi dan dibuat",
		"customToken": customToken,
		"email":       body.Email,
		"uid":         userRecord.UID,
		"role":        roleName,
	})
}
