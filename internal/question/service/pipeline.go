package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// pipelineTimeout is the maximum time Stage 3 (LLM) may take before the
// pipeline falls back to the Stage 2 structural score.
const pipelineTimeout = 10 * time.Second

// PipelineEvaluator runs the three-stage evaluation in sequence:
//
//	Stage 1 — keyword filter: fast rejection for clearly incomplete answers.
//	Stage 2 — structural coverage: deterministic per-concept match report.
//	Stage 3 — LLM judgment: semantic scoring with Stage 2 as grounding context.
//
// Stage 3 is optional. When llm is nil (no API key configured), the pipeline
// stops at Stage 2. If Stage 3 fails or times out, it also falls back to Stage 2
// so the user always gets a result.
type PipelineEvaluator struct {
	llm LLMEvaluator // nil when not configured
}

// NewPipelineEvaluator creates a PipelineEvaluator.
//
// Pass a nil LLMEvaluator to run Stages 1+2 only. This is the correct way to
// express "no LLM configured" — do not pass a non-nil concrete pointer wrapped
// in the interface, because that passes a nil-check silently.
//
//	client := NewGeminiClient(os.Getenv("GEMINI_API_KEY"))
//	var llm LLMEvaluator
//	if client != nil {
//	    llm = client  // only assign when the concrete pointer is non-nil
//	}
//	evaluator := NewPipelineEvaluator(llm)
func NewPipelineEvaluator(llm LLMEvaluator) *PipelineEvaluator {
	return &PipelineEvaluator{llm: llm}
}

// Evaluate satisfies the Evaluator interface.
// It delegates to EvaluateAnswer for Stages 1+2 and optionally calls the LLM
// for Stage 3 when configured.
func (p *PipelineEvaluator) Evaluate(answer, rawGuide string) EvaluationResult {
	// Run Stage 1 + Stage 2 via the shared structural path.
	// This call is always fast (microseconds to low milliseconds).
	structural := EvaluateAnswer(answer, rawGuide)

	// If Stage 1 rejected the answer, stop here — do not call the LLM.
	if structural.EvaluatedBy == "filter_rejected" {
		return structural
	}

	// If no LLM is configured, return the structural result directly.
	if p.llm == nil {
		return structural
	}

	// Stage 3: LLM judgment.
	// Build the coverage report again so it can be passed as grounding context.
	answerLower := strings.ToLower(strings.TrimSpace(answer))
	guide := ParseAnswerGuide(rawGuide)
	report := buildCoverageReport(answerLower, guide)

	llmResult, err := p.callLLM(context.Background(), answer, rawGuide, report)
	if err != nil {
		// Stage 3 failed or timed out — fall back to Stage 2 and log.
		log.Printf("pipeline: stage 3 fallback (structural score=%d): %v", structural.StructuralScore, err)
		return structural
	}

	return llmResult
}

// callLLM wraps the LLM call with a timeout and up to two retries on validation failure.
func (p *PipelineEvaluator) callLLM(ctx context.Context, answer, _ string, report CoverageReport) (EvaluationResult, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, pipelineTimeout)
	defer cancel()

	req := LLMRequest{
		CandidateAnswer: answer,
		Coverage:        report,
	}

	const maxRetries = 2
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err := p.llm.Evaluate(timeoutCtx, req)
		if err != nil {
			lastErr = err
			// Context deadline exceeded — no point retrying.
			if timeoutCtx.Err() != nil {
				return EvaluationResult{}, fmt.Errorf("llm timeout after %v: %w", pipelineTimeout, timeoutCtx.Err())
			}
			log.Printf("pipeline: llm attempt %d/%d failed: %v", attempt, maxRetries, err)
			continue
		}

		if err := validateLLMResponse(resp); err != nil {
			lastErr = err
			log.Printf("pipeline: llm attempt %d/%d invalid response: %v", attempt, maxRetries, err)
			continue
		}

		// Valid response — convert to EvaluationResult.
		return EvaluationResult{
			Score:           resp.Score,
			Correct:         resp.Correct,
			Strengths:       resp.Strengths,
			Gaps:            resp.Gaps,
			Summary:         resp.Summary,
			StructuralScore: report.StructuralScore,
			EvaluatedBy:     "llm",
		}, nil
	}

	return EvaluationResult{}, fmt.Errorf("llm all retries exhausted: %w", lastErr)
}

// validateLLMResponse checks that all required fields are present and in range.
func validateLLMResponse(r LLMResponse) error {
	if r.Score < 0 || r.Score > 100 {
		return fmt.Errorf("score %d out of range [0,100]", r.Score)
	}
	if r.Strengths == nil {
		return fmt.Errorf("strengths is nil")
	}
	if r.Gaps == nil {
		return fmt.Errorf("gaps is nil")
	}
	if len(strings.TrimSpace(r.Summary)) == 0 {
		return fmt.Errorf("summary is empty")
	}
	return nil
}
