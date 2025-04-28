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
		<h2>Verifikasi Email</h2>
		<p>Klik tombol di bawah ini untuk memverifikasi akun kamu:</p>
		<a href="%s" style="display:inline-block;padding:10px 20px;background:#2563eb;color:white;text-decoration:none;border-radius:6px;">Verifikasi Akun</a>
		<p>Link ini akan kadaluarsa dalam 15 menit.</p>
	`, link)

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
