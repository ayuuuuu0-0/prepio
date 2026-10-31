package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// Mailer defines the email sending capability for Prepio auth.
type Mailer interface {
	SendOTP(ctx context.Context, toEmail, code string) error
	IsDev() bool
}

// NewMailerFromEnv creates the appropriate Mailer implementation based on environment variables.
// Priority:
// 1. SMTP if SMTP_USER and SMTP_PASSWORD are set (e.g. Gmail SMTP with App Password)
// 2. Resend if RESEND_API_KEY is set
// 3. Fallback to Dev/Mock mode (logs code to console, dev_code returned in dev)
func NewMailerFromEnv() Mailer {
	smtpUser := strings.TrimSpace(os.Getenv("SMTP_USER"))
	smtpPass := strings.TrimSpace(os.Getenv("SMTP_PASSWORD"))
	if smtpUser != "" && smtpPass != "" {
		host := os.Getenv("SMTP_HOST")
		if host == "" {
			host = "smtp.gmail.com"
		}
		port := os.Getenv("SMTP_PORT")
		if port == "" {
			port = "587"
		}
		fromName := os.Getenv("SMTP_FROM_NAME")
		if fromName == "" {
			fromName = "Prepio"
		}
		log.Printf("mailer: using SMTP (%s:%s as %s)", host, port, smtpUser)
		return NewSMTPMailer(host, port, smtpUser, smtpPass, fromName)
	}

	resendKey := strings.TrimSpace(os.Getenv("RESEND_API_KEY"))
	if resendKey != "" && resendKey != "mock" && resendKey != "dev" {
		fromEmail := os.Getenv("RESEND_FROM_EMAIL")
		log.Printf("mailer: using Resend API")
		return NewResendMailer(resendKey, fromEmail)
	}

	log.Printf("mailer: using Dev/Mock mode (no SMTP or Resend credentials configured)")
	return NewResendMailer("", "")
}

// -----------------------------------------------------------------------------
// SMTP Mailer (e.g. Gmail SMTP with App Password)
// -----------------------------------------------------------------------------

// SMTPMailer sends emails via standard SMTP with STARTTLS (port 587).
type SMTPMailer struct {
	host     string
	port     string
	user     string
	password string
	fromName string
}

// NewSMTPMailer creates a new SMTPMailer.
func NewSMTPMailer(host, port, user, password, fromName string) *SMTPMailer {
	return &SMTPMailer{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		fromName: fromName,
	}
}

func (m *SMTPMailer) IsDev() bool {
	return m.user == "" || m.password == ""
}

// SendOTP sends the verification email using SMTP.
func (m *SMTPMailer) SendOTP(ctx context.Context, toEmail, code string) error {
	if m.IsDev() {
		log.Printf("[OTP DEV MODE] Verification code for %s is: %s", toEmail, code)
		return nil
	}

	fromHeader := fmt.Sprintf("%s <%s>", m.fromName, m.user)
	subject := fmt.Sprintf("%s is your Prepio verification code", code)

	msg := buildMIMEMessage(fromHeader, toEmail, subject, renderOTPEmail(code))
	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	auth := smtp.PlainAuth("", m.user, m.password, m.host)

	errChan := make(chan error, 1)
	go func() {
		errChan <- smtp.SendMail(addr, auth, m.user, []string{toEmail}, msg)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("send email via smtp: %w", err)
		}
		log.Printf("mailer: successfully sent OTP email to %s via %s", toEmail, m.host)
		return nil
	}
}

// -----------------------------------------------------------------------------
// Resend Mailer (REST API)
// -----------------------------------------------------------------------------

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

// SendOTP delivers the 6-digit verification code via Resend.
func (m *ResendMailer) SendOTP(ctx context.Context, toEmail, code string) error {
	if m.IsDev() {
		log.Printf("[OTP DEV MODE] Verification code for %s is: %s", toEmail, code)
		return nil
	}

	payload := resendPayload{
		From:    m.fromEmail,
		To:      []string{toEmail},
		Subject: fmt.Sprintf("%s is your Prepio verification code", code),
		HTML:    renderOTPEmail(code),
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

	log.Printf("mailer: successfully sent OTP email to %s via Resend", toEmail)
	return nil
}

// -----------------------------------------------------------------------------
// Mock Mailer (for tests)
// -----------------------------------------------------------------------------

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

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func buildMIMEMessage(from, to, subject, htmlBody string) []byte {
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)
	return msg.Bytes()
}

func renderOTPEmail(code string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Your Prepio Verification Code</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #0E1015; margin: 0; padding: 40px; color: #E8EAED;">
  <div style="max-width: 480px; margin: 0 auto; background-color: #13161F; border: 1px solid #2E3347; border-radius: 16px; padding: 32px; text-align: center;">
    <h1 style="color: #7C6EF5; font-size: 24px; font-weight: bold; margin-top: 0; margin-bottom: 8px; letter-spacing: 1px;">PREPIO</h1>
    <p style="color: #8C93A8; font-size: 14px; margin-top: 0; margin-bottom: 24px;">Career RPG Verification Code</p>
    
    <p style="font-size: 15px; color: #E8EAED; margin-bottom: 20px;">Use the 6-digit code below to verify your email and activate your account:</p>
    
    <div style="background-color: #1A1D27; border: 1px solid #3E445B; border-radius: 12px; padding: 18px 24px; margin: 24px 0; display: inline-block;">
      <span style="font-size: 34px; font-weight: 800; letter-spacing: 8px; color: #F5A623; font-family: monospace;">%s</span>
    </div>
    
    <p style="color: #8C93A8; font-size: 13px; line-height: 1.5; margin-top: 24px; margin-bottom: 0;">
      This code is valid for <strong>5 minutes</strong>. If you did not request this verification code, you can safely ignore this email.
    </p>
  </div>
</body>
</html>`, code)
}
