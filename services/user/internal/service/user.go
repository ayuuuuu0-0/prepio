package service

import (
	"context"

	"github.com/prepio/prepio/services/user/internal/dto"
	"github.com/prepio/prepio/services/user/internal/store"
)

// UserService handles profile and device operations.
type UserService struct {
	users   *store.UserStore
	devices *store.UserDeviceStore
}

// NewUserService creates a UserService.
func NewUserService(users *store.UserStore, devices *store.UserDeviceStore) *UserService {
	return &UserService{users: users, devices: devices}
}

// GetProfile returns the authenticated user's profile.
func (s *UserService) GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	resp := toUserResponse(user)
	return &resp, nil
}

// UpdateProfile updates the authenticated user's profile.
func (s *UserService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	timezone := user.Timezone
	if req.Timezone != nil && len(*req.Timezone) > 0 {
		timezone = *req.Timezone
	}

	reminderTime := user.ReminderTime
	if req.ReminderTime != nil && len(*req.ReminderTime) > 0 {
		reminderTime = *req.ReminderTime
	}

	updated, err := s.users.UpdateProfile(ctx, userID, timezone, reminderTime, user.ActiveCharID)
	if err != nil {
		return nil, err
	}

	resp := toUserResponse(updated)
	return &resp, nil
}

// RegisterDevice upserts an FCM token for the user.
func (s *UserService) RegisterDevice(ctx context.Context, userID string, req dto.RegisterDeviceRequest) (*dto.DeviceResponse, error) {
	if len(req.FCMToken) == 0 || len(req.Platform) == 0 {
		return nil, ErrInvalidRequest
	}

	device, err := s.devices.Upsert(ctx, userID, req.FCMToken, req.Platform)
	if err != nil {
		return nil, err
	}

	return &dto.DeviceResponse{
		ID:       device.ID,
		Platform: device.Platform,
	}, nil
}

// DeleteDevice removes a registered device.
func (s *UserService) DeleteDevice(ctx context.Context, userID, deviceID string) error {
	if len(deviceID) == 0 {
		return ErrInvalidRequest
	}

	deleted, err := s.devices.Delete(ctx, userID, deviceID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrDeviceNotFound
	}
	return nil
}

func toUserResponse(user *store.User) dto.UserResponse {
	return dto.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		Timezone:     user.Timezone,
		ActiveCharID: user.ActiveCharID,
		ReminderTime: user.ReminderTime,
	}
}
