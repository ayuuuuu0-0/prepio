package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Mailer defines the email sending capability for Prepio auth.
type Mailer interface {
	SendOTP(ctx context.Context, toEmail, code string) error
	IsDev() bool
}

// ResendMailer sends emails via Resend's REST API.
type ResendMailer struct {
	apiKey    string
	fromEmail string
	client    *http.Client
}

// NewResendMailer creates a new ResendMailer.
func NewResendMailer(apiKey, fromEmail string) *ResendMailer {
	if fromEmail == "" {
		fromEmail = "Prepio <onboarding@resend.dev>"
	}
	return &ResendMailer{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *ResendMailer) IsDev() bool {
	return m.apiKey == "" || m.apiKey == "mock" || m.apiKey == "dev"
}

type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

// SendOTP delivers the 6-digit verification code to the recipient email.
func (m *ResendMailer) SendOTP(ctx context.Context, toEmail, code string) error {
	// Fallback/development mode: if API key is not configured or in tests, log the code
	if m.IsDev() {
		log.Printf("[OTP DEV MODE] Verification code for %s is: %s", toEmail, code)
		return nil
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Your Prepio Verification Code</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #0E1015; margin: 0; padding: 40px; color: #E8EAED;">
  <div style="max-width: 480px; margin: 0 auto; background-color: #13161F; border: 1px solid #2E3347; border-radius: 16px; padding: 32px; text-align: center;">
    <h1 style="color: #7C6EF5; font-size: 24px; font-weight: bold; margin-top: 0; margin-bottom: 8px;">PREPIO</h1>
    <p style="color: #8C93A8; font-size: 14px; margin-top: 0; margin-bottom: 24px;">Career RPG Verification Code</p>
    
    <p style="font-size: 15px; color: #E8EAED; margin-bottom: 20px;">Use the 6-digit code below to sign in or create your account:</p>
    
    <div style="background-color: #1A1D27; border: 1px solid #3E445B; border-radius: 12px; padding: 18px 24px; margin: 24px 0; display: inline-block;">
      <span style="font-size: 34px; font-weight: 800; letter-spacing: 8px; color: #F5A623; font-family: monospace;">%s</span>
    </div>
    
    <p style="color: #8C93A8; font-size: 13px; line-height: 1.5; margin-top: 24px; margin-bottom: 0;">
      This code is valid for <strong>5 minutes</strong>. If you did not request this code, you can safely ignore this email.
    </p>
  </div>
</body>
</html>`, code)

	payload := resendPayload{
		From:    m.fromEmail,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("%s is your Prepio verification code", code),
		HTML:    html,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal resend payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create resend request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", m.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("send email via resend: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("resend api error: status %d", resp.StatusCode)
	}

	return nil
}

// MockMailer records sent OTPs for unit and smoke testing.
type MockMailer struct {
	LastEmail string
	LastCode  string
}

func (m *MockMailer) SendOTP(ctx context.Context, toEmail, code string) error {
	m.LastEmail = toEmail
	m.LastCode = code
	return nil
}

func (m *MockMailer) IsDev() bool {
	return true
}
