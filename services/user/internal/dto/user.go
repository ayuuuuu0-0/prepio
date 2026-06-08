package dto

// UserResponse is the public user profile shape.
type UserResponse struct {
	ID           string  `json:"id"`
	Email        string  `json:"email"`
	Username     string  `json:"username"`
	Timezone     string  `json:"timezone,omitempty"`
	ActiveCharID *string `json:"active_char_id,omitempty"`
	ReminderTime string  `json:"reminder_time,omitempty"`
}

// UpdateProfileRequest is the body for PATCH /api/v1/users/me.
type UpdateProfileRequest struct {
	Timezone     *string `json:"timezone"`
	ReminderTime *string `json:"reminder_time"`
}

// RegisterDeviceRequest is the body for POST /api/v1/users/me/devices.
type RegisterDeviceRequest struct {
	FCMToken string `json:"fcm_token"`
	Platform string `json:"platform"`
}

// DeviceResponse is the public device shape.
type DeviceResponse struct {
	ID       string `json:"id"`
	Platform string `json:"platform"`
}
