package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/winterheatherica/tokoaku-backend/config"
)

type EmailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Html    string `json:"html"`
}

func SendVerificationEmail(to string, token string) error {
	link := fmt.Sprintf(
		"%s/verify?email=%s&token=%s",
		config.App.FrontendBaseURL,
		to,
		token,
	)

	html := fmt.Sprintf(`
	<div style="font-family:'Segoe UI',sans-serif;max-width:480px;margin:0 auto;padding:24px;background-color:#f5fffd;border:1px solid #8FCFBC;border-radius:12px;box-shadow:0 4px 12px rgba(58,175,169,0.15);color:#22372B;">
		<h2 style="color:#3AAFA9;margin-top:0;">Halo dari %s 👋</h2>
		<p style="font-size:1rem;line-height:1.6;">
		Terima kasih telah mendaftar! Kami hanya butuh satu langkah lagi untuk mengaktifkan akun kamu.
		</p>
		<div style="text-align:center;margin:28px 0;">
		<a href="%s" style="display:inline-block;padding:12px 24px;background-color:#3AAFA9;color:white;font-weight:600;text-decoration:none;border-radius:8px;font-size:1rem;box-shadow:0 4px 10px rgba(58,175,169,0.25);">
			✔ Verifikasi Sekarang
		</a>
		</div>
		<p style="font-size:0.95rem;color:#555;">
		Tombol di atas akan membawa kamu ke halaman verifikasi. Link ini akan kadaluarsa dalam <strong>15 menit</strong>.
		</p>
		<hr style="border:none;border-top:1px solid #8FCFBC;margin:24px 0;" />
		<p style="font-size:0.85rem;color:#888;">
		Jika kamu tidak merasa mendaftar ke %s, kamu bisa mengabaikan email ini.
		</p>
		<p style="font-size:0.85rem;color:#8FCFBC;text-align:center;margin-top:32px;">✨ Buat pengalaman berbelanja jadi lebih menyenangkan bersama %s ✨</p>
	</div>
	`, config.App.PlatformName, link, config.App.PlatformName, config.App.PlatformName)

	payload := EmailPayload{
		From:    config.Email.Sender,
		To:      to,
		Subject: "Verifikasi Email " + config.App.PlatformName,
		Html:    html,
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+config.Email.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("failed to send email, status: %d", res.StatusCode)
	}

	return nil
}
