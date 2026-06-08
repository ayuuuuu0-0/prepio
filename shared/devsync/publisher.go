package devsync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prepio/prepio/shared/events"
)

// Publisher forwards domain events to internal HTTP endpoints for local dev without Kafka.
type Publisher struct {
	client          *http.Client
	streakURL       string
	progressURL     string
	notificationURL string
}

// Config holds service base URLs for dev event forwarding.
type Config struct {
	StreakURL       string
	ProgressURL     string
	NotificationURL string
}

// NewPublisher creates a dev-sync event publisher.
func NewPublisher(cfg Config) *Publisher {
	return &Publisher{
		client:          &http.Client{Timeout: 10 * time.Second},
		streakURL:       cfg.StreakURL,
		progressURL:     cfg.ProgressURL,
		notificationURL: cfg.NotificationURL,
	}
}

// NewPublisherFromEnv builds a Publisher using standard service URL environment variables.
func NewPublisherFromEnv() *Publisher {
	return NewPublisher(Config{
		StreakURL:       envOrDefault("STREAK_SERVICE_URL", "http://localhost:8083"),
		ProgressURL:     envOrDefault("PROGRESS_SERVICE_URL", "http://localhost:8084"),
		NotificationURL: envOrDefault("NOTIFICATION_SERVICE_URL", "http://localhost:8085"),
	})
}

// Publish routes the event to the appropriate internal HTTP handlers.
func (p *Publisher) Publish(ctx context.Context, topic, key string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	switch topic {
	case events.TopicQuestionAnswered:
		if err := p.post(ctx, p.streakURL+"/internal/events/question-answered", body); err != nil {
			return err
		}
		return p.post(ctx, p.progressURL+"/internal/events/question-answered", body)
	case events.TopicStreakUpdated:
		if err := p.post(ctx, p.progressURL+"/internal/events/streak-updated", body); err != nil {
			return err
		}
		return p.post(ctx, p.notificationURL+"/internal/events/streak-updated", body)
	case events.TopicProgressUpdated:
		return p.post(ctx, p.notificationURL+"/internal/events/progress-updated", body)
	default:
		return nil
	}
}

// Close is a no-op for the HTTP publisher.
func (p *Publisher) Close() error {
	return nil
}

func (p *Publisher) post(ctx context.Context, url string, body []byte) error {
	if len(url) == 0 {
		return fmt.Errorf("dev sync url is required")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("post %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("post %s: status %d", url, resp.StatusCode)
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); len(v) > 0 {
		return v
	}
	return fallback
}
