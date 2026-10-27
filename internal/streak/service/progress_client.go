package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ProgressGemClient calls the progress service internal gem API.
type ProgressGemClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewProgressGemClient creates a ProgressGemClient.
func NewProgressGemClient(baseURL string) *ProgressGemClient {
	return &ProgressGemClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// DeductGems calls POST /internal/progress/{userID}/gems/deduct.
func (c *ProgressGemClient) DeductGems(ctx context.Context, userID string, amount int, reason string) error {
	body, err := json.Marshal(map[string]any{
		"amount": amount,
		"reason": reason,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/internal/progress/%s/gems/deduct", c.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return ErrInsufficientGems
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("deduct gems status %d", resp.StatusCode)
	}
	return nil
}
