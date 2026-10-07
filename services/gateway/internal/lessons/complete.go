// Package lessons aggregates lesson completion with the rewards Progress computed.
// Aggregation does not change ownership: Journey grades and completes, Progress rewards.
package lessons

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

const (
	rewardsAttempts = 5
	rewardsBackoff  = 120 * time.Millisecond
)

// Handler completes an attempt and attaches Progress rewards to the summary.
type Handler struct {
	questionURL string
	progressURL string
	client      *http.Client
}

// NewHandler creates a Handler.
func NewHandler(questionURL, progressURL string) *Handler {
	return &Handler{questionURL: questionURL, progressURL: progressURL, client: &http.Client{Timeout: 10 * time.Second}}
}

// Complete handles POST /api/v1/attempts/{id}/complete.
func (h *Handler) Complete(w http.ResponseWriter, r *http.Request) {
	token := middleware.ExtractBearerToken(r)
	if len(token) == 0 {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}
	attemptID := url.PathEscape(r.PathValue("id"))

	status, body, err := h.do(r.Context(), http.MethodPost, h.questionURL+"/api/v1/attempts/"+attemptID+"/complete", token)
	if err != nil {
		response.Error(w, http.StatusBadGateway, constants.ErrInternal, "lesson service unavailable")
		return
	}
	if status != http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(body)
		return
	}

	var envelope struct {
		Data map[string]json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil || envelope.Data == nil {
		response.Error(w, http.StatusBadGateway, constants.ErrInternal, "invalid lesson response")
		return
	}

	rewards := h.fetchRewards(r.Context(), attemptID, token)
	if rewards != nil {
		envelope.Data["rewards"] = rewards
		envelope.Data["rewards_pending"] = json.RawMessage("false")
	} else {
		envelope.Data["rewards"] = json.RawMessage("null")
		envelope.Data["rewards_pending"] = json.RawMessage("true")
	}
	response.Data(w, http.StatusOK, envelope.Data)
}

// fetchRewards reads the attempt's rewards from Progress. Progress applies the
// completion event asynchronously, so it retries briefly before reporting pending.
func (h *Handler) fetchRewards(ctx context.Context, attemptID, token string) json.RawMessage {
	for i := 0; i < rewardsAttempts; i++ {
		status, body, err := h.do(ctx, http.MethodGet, h.progressURL+"/api/v1/progress/attempts/"+attemptID+"/rewards", token)
		if err == nil && status == http.StatusOK {
			var envelope struct {
				Data json.RawMessage `json:"data"`
			}
			if json.Unmarshal(body, &envelope) == nil && len(envelope.Data) > 0 {
				return envelope.Data
			}
			return nil
		}
		if err == nil && status != http.StatusNotFound {
			return nil
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(rewardsBackoff):
		}
	}
	return nil
}

func (h *Handler) do(ctx context.Context, method, target, token string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, target, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := h.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}
	return res.StatusCode, body, nil
}
