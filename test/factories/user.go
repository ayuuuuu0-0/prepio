package factories

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/prepio/prepio/constants"
)

// UserFactory builds test user data with sensible defaults.
type UserFactory struct {
	Email    string
	Username string
	Password string
	Timezone string
}

// NewUserFactory returns a UserFactory with unique email and username.
func NewUserFactory() UserFactory {
	suffix := uuid.NewString()[:8]
	return UserFactory{
		Email:    fmt.Sprintf("user-%s@example.test", suffix),
		Username: fmt.Sprintf("user_%s", suffix),
		Password: "password123",
		Timezone: constants.DefaultTimezone,
	}
}

// WithEmail overrides the email.
func (f UserFactory) WithEmail(email string) UserFactory {
	f.Email = email
	return f
}

// WithUsername overrides the username.
func (f UserFactory) WithUsername(username string) UserFactory {
	f.Username = username
	return f
}

// WithPassword overrides the password.
func (f UserFactory) WithPassword(password string) UserFactory {
	f.Password = password
	return f
}
