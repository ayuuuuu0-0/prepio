package service

import (
	"context"
	"fmt"
)

// LLMEvaluator is the interface for an AI-powered answer evaluator.
// It receives the question text, the candidate's answer, and the Stage 2
// coverage report, and returns a scored result.
//
// The interface is intentionally minimal. Wire a concrete implementation
// only when an API key is configured; pass nil to NewPipelineEvaluator when
// it isn't, and the pipeline will degrade gracefully to Stage 2 only.
type LLMEvaluator interface {
	Evaluate(ctx context.Context, req LLMRequest) (LLMResponse, error)
}

// LLMRequest is the input handed to the LLM stage.
type LLMRequest struct {
	// QuestionBody is the full question text shown to the candidate.
	QuestionBody string

	// CandidateAnswer is the raw submission from the user.
	CandidateAnswer string

	// Coverage is the Stage 2 structural report used as grounding context.
	// The LLM is instructed to trust its own judgment about paraphrases, not
	// to treat the "missed" list as definitively wrong.
	Coverage CoverageReport
}

// LLMResponse is the parsed output expected from the LLM.
// All fields are required; the pipeline validates each one before trusting it.
type LLMResponse struct {
	Score     int      `json:"score"`
	Correct   bool     `json:"correct"`
	Strengths []string `json:"strengths"`
	Gaps      []string `json:"gaps"`
	Summary   string   `json:"summary"`
}

// --- Gemini placeholder -------------------------------------------------------

// GeminiClient will become a live Gemini API client once an API key is supplied.
// For now it returns a clear error so the pipeline falls back to Stage 2 and
// logs the fact that Stage 3 is not yet active.
//
// Replace the body of Evaluate with the real Gemini REST call when the key is
// available. The interface, request/response shapes, and pipeline wiring do not
// need to change.
type GeminiClient struct {
	// APIKey will be the Gemini API key. Currently unused.
	APIKey string
}

// NewGeminiClient returns a GeminiClient if an API key is present, or nil if not.
// The caller must check for nil before assigning to an LLMEvaluator variable to
// avoid the Go nil-interface trap (a non-nil interface wrapping a nil pointer
// passes a != nil check silently).
func NewGeminiClient(apiKey string) *GeminiClient {
	if len(apiKey) == 0 {
		return nil
	}
	return &GeminiClient{APIKey: apiKey}
}

// Evaluate is the placeholder implementation. It always returns an error so the
// pipeline falls back to the Stage 2 structural score.
// TODO: replace this body with the real Gemini generateContent call.
func (g *GeminiClient) Evaluate(_ context.Context, _ LLMRequest) (LLMResponse, error) {
	return LLMResponse{}, fmt.Errorf("gemini client: not yet implemented — falling back to structural score")
}
