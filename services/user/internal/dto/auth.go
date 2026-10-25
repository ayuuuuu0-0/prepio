package dto

// RegisterRequest is the body for POST /api/v1/auth/register.
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Timezone string `json:"timezone"`
}

// LoginRequest is the body for POST /api/v1/auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SendOTPRequest is the body for POST /api/v1/auth/otp/send.
type SendOTPRequest struct {
	Email string `json:"email"`
}

// SendOTPResponse is returned by POST /api/v1/auth/otp/send.
type SendOTPResponse struct {
	Message   string `json:"message"`
	ExpiresIn int    `json:"expires_in"` // in seconds
}

// VerifyOTPRequest is the body for POST /api/v1/auth/otp/verify.
type VerifyOTPRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Username string `json:"username,omitempty"`
}

// RefreshRequest is the body for POST /api/v1/auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthResponse is returned by register, login, and verify OTP.
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}
