package service

import (
	"encoding/json"
	"strings"
)

// Concept describes a single expected idea in an answer guide.
type Concept struct {
	// Name is the primary label for this concept.
	Name string `json:"name"`

	// Aliases are alternative phrasings that count as a match (e.g. "hashmap" for "hash map").
	Aliases []string `json:"aliases,omitempty"`

	// Required marks this concept as mandatory for the Stage 1 keyword filter.
	// Optional concepts still contribute to the structural score in Stage 2.
	Required bool `json:"required"`

	// Requires lists concept names that are logical prerequisites of this one.
	// Used as grounding context in the Stage 3 LLM prompt; not enforced as a hard block.
	Requires []string `json:"requires,omitempty"`
}

// AnswerGuide is the structured rubric attached to a question.
type AnswerGuide struct {
	Concepts []Concept `json:"concepts"`
}

// ParseAnswerGuide parses an answer guide string into structured form.
// It tries JSON first, then falls back to the legacy "concepts: a|b|c" pipe format.
// This lets existing questions work without any data migration.
func ParseAnswerGuide(raw string) AnswerGuide {
	raw = strings.TrimSpace(raw)
	if len(raw) == 0 {
		return AnswerGuide{}
	}

	// Try structured JSON first (new format starts with '{').
	if strings.HasPrefix(raw, "{") {
		var guide AnswerGuide
		if err := json.Unmarshal([]byte(raw), &guide); err == nil && len(guide.Concepts) > 0 {
			return guide
		}
	}

	// Fall back to legacy pipe-separated format.
	return parseLegacyGuide(raw)
}

// parseLegacyGuide handles the old "concepts: a|b|c" format.
// All extracted concepts are treated as required (matching old behaviour).
func parseLegacyGuide(raw string) AnswerGuide {
	lower := strings.ToLower(raw)
	const prefix = "concepts:"
	if !strings.HasPrefix(lower, prefix) {
		return AnswerGuide{}
	}

	rest := strings.TrimSpace(raw[len(prefix):])
	// Only parse the first line; the rest is free-text description.
	if idx := strings.Index(rest, "\n"); idx >= 0 {
		rest = rest[:idx]
	}

	parts := strings.Split(rest, "|")
	concepts := make([]Concept, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			concepts = append(concepts, Concept{
				Name:     part,
				Required: true,
			})
		}
	}
	return AnswerGuide{Concepts: concepts}
}

// RequiredConcepts returns only the concepts marked as required.
func (g AnswerGuide) RequiredConcepts() []Concept {
	out := make([]Concept, 0, len(g.Concepts))
	for _, c := range g.Concepts {
		if c.Required {
			out = append(out, c)
		}
	}
	return out
}
