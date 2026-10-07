// Package lesson defines the authored lesson model (content-as-code), its typed
// step payloads, the answer-free client projection, and content validation.
package lesson

import (
	"encoding/json"
	"fmt"
)

// Step types.
const (
	StepIntro     = "intro"
	StepMCQ       = "mcq"
	StepTrueFalse = "true_false"
	StepFillBlank = "fill_blank"
	StepArrange   = "arrange"
	StepProse     = "prose"
)

// Lesson statuses. Only StatusPublished is served.
const (
	StatusDraft      = "draft"
	StatusReview     = "review"
	StatusPublished  = "published"
	StatusDeprecated = "deprecated"
)

// Lesson kinds.
const (
	KindLesson = "lesson"
	KindBoss   = "boss"
)

// Node types.
const (
	NodeTypeLesson = "lesson"
	NodeTypeBoss   = "boss"
)

// IsGraded reports whether a step type is graded.
func IsGraded(stepType string) bool {
	return stepType != StepIntro
}

// Lesson is one authored lesson file under content/lessons.
type Lesson struct {
	Slug       string     `yaml:"slug" json:"slug"`
	Title      string     `yaml:"title" json:"title"`
	Node       string     `yaml:"node" json:"node"`
	Kind       string     `yaml:"kind" json:"kind"`
	Status     string     `yaml:"status" json:"-"`
	Difficulty string     `yaml:"difficulty" json:"difficulty"`
	EstMinutes int        `yaml:"est_minutes" json:"est_minutes"`
	Skills     []SkillRef `yaml:"skills" json:"skills"`
	Summary    []string   `yaml:"summary" json:"summary"`
	Steps      []Step     `yaml:"steps" json:"steps"`
}

// SkillRef maps a lesson to a skill with a weight.
type SkillRef struct {
	Skill  string  `yaml:"skill" json:"skill"`
	Weight float64 `yaml:"weight" json:"weight"`
}

// Step is the authoring shape of one lesson step. Which fields apply depends on Type.
type Step struct {
	ID   string `yaml:"id" json:"id"`
	Type string `yaml:"type" json:"type"`

	// intro
	Beats    []Beat `yaml:"beats,omitempty" json:"beats,omitempty"`
	MediaURL string `yaml:"mediaUrl,omitempty" json:"media_url,omitempty"`

	// mcq, arrange, prose
	Prompt string `yaml:"prompt,omitempty" json:"prompt,omitempty"`
	// true_false
	Statement string `yaml:"statement,omitempty" json:"statement,omitempty"`

	// mcq
	Options      []string      `yaml:"options,omitempty" json:"options,omitempty"`
	Explanations *Explanations `yaml:"explanations,omitempty" json:"explanations,omitempty"`

	// mcq (int index) and true_false (bool)
	Answer any `yaml:"answer,omitempty" json:"answer,omitempty"`

	// fill_blank
	Code    string   `yaml:"code,omitempty" json:"code,omitempty"`
	Bank    []string `yaml:"bank,omitempty" json:"bank,omitempty"`
	Answers []string `yaml:"answers,omitempty" json:"answers,omitempty"`

	// arrange
	Items []string `yaml:"items,omitempty" json:"items,omitempty"`
	Order []int    `yaml:"order,omitempty" json:"order,omitempty"`

	// prose
	Rubric   *Rubric `yaml:"rubric,omitempty" json:"rubric,omitempty"`
	MinChars int     `yaml:"minChars,omitempty" json:"min_chars,omitempty"`

	// true_false, fill_blank, arrange, prose
	Explanation string `yaml:"explanation,omitempty" json:"explanation,omitempty"`
}

// Beat is one beat of an intro.
type Beat struct {
	Text       string `yaml:"text" json:"text"`
	Emphasis   string `yaml:"emphasis,omitempty" json:"emphasis,omitempty"`
	Visual     string `yaml:"visual,omitempty" json:"visual,omitempty"`
	DurationMs int    `yaml:"durationMs" json:"duration_ms"`
}

// Explanations holds the mcq explanation for the right option and each wrong option.
type Explanations struct {
	Correct string         `yaml:"correct" json:"correct"`
	Wrong   map[int]string `yaml:"wrong" json:"wrong"`
}

// Rubric is the concept list the prose evaluator scores against.
type Rubric struct {
	Concepts []RubricConcept `yaml:"concepts" json:"concepts"`
}

// RubricConcept is one expected idea in a prose answer.
type RubricConcept struct {
	Name     string   `yaml:"name" json:"name"`
	Aliases  []string `yaml:"aliases,omitempty" json:"aliases,omitempty"`
	Required bool     `yaml:"required" json:"required"`
}

// World is one authored world file under content/worlds.
type World struct {
	Slug        string `yaml:"slug" json:"slug"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Theme       string `yaml:"theme" json:"theme"`
	Order       int    `yaml:"order" json:"order"`
	Nodes       []Node `yaml:"nodes" json:"nodes"`
}

// Node is one milestone in a world. It binds exactly one lesson.
type Node struct {
	Slug     string   `yaml:"slug" json:"slug"`
	Label    string   `yaml:"label" json:"label"`
	Type     string   `yaml:"type" json:"type"`
	Requires []string `yaml:"requires,omitempty" json:"requires,omitempty"`
}

// Content is everything loaded from content/.
type Content struct {
	Worlds  []World
	Lessons []Lesson
}

// Typed, stored step payloads. Authoring Step values are converted to these
// before they are written to lesson_steps.payload.

// IntroPayload is the stored payload of an intro step.
type IntroPayload struct {
	Beats    []Beat `json:"beats"`
	MediaURL string `json:"media_url,omitempty"`
}

// MCQPayload is the stored payload of a multiple-choice step.
type MCQPayload struct {
	Prompt       string       `json:"prompt"`
	Options      []string     `json:"options"`
	Answer       int          `json:"answer"`
	Explanations Explanations `json:"explanations"`
}

// TrueFalsePayload is the stored payload of a true/false step.
type TrueFalsePayload struct {
	Statement   string `json:"statement"`
	Answer      bool   `json:"answer"`
	Explanation string `json:"explanation"`
}

// FillBlankPayload is the stored payload of a fill-in-the-blank step.
type FillBlankPayload struct {
	Code        string   `json:"code"`
	Bank        []string `json:"bank"`
	Answers     []string `json:"answers"`
	Explanation string   `json:"explanation"`
}

// ArrangePayload is the stored payload of an arrange-in-order step.
type ArrangePayload struct {
	Prompt      string   `json:"prompt"`
	Items       []string `json:"items"`
	Order       []int    `json:"order"`
	Explanation string   `json:"explanation"`
}

// ProsePayload is the stored payload of a rubric-graded prose step.
type ProsePayload struct {
	Prompt      string `json:"prompt"`
	Rubric      Rubric `json:"rubric"`
	MinChars    int    `json:"min_chars"`
	Explanation string `json:"explanation"`
}

// StoredPayload converts an authored step to its typed payload. It assumes the
// step has passed validation.
func (s Step) StoredPayload() (any, error) {
	switch s.Type {
	case StepIntro:
		return IntroPayload{Beats: s.Beats, MediaURL: s.MediaURL}, nil
	case StepMCQ:
		answer, ok := asInt(s.Answer)
		if !ok || s.Explanations == nil {
			return nil, fmt.Errorf("step %q: invalid mcq", s.ID)
		}
		return MCQPayload{Prompt: s.Prompt, Options: s.Options, Answer: answer, Explanations: *s.Explanations}, nil
	case StepTrueFalse:
		answer, ok := s.Answer.(bool)
		if !ok {
			return nil, fmt.Errorf("step %q: true_false answer must be a boolean", s.ID)
		}
		return TrueFalsePayload{Statement: s.Statement, Answer: answer, Explanation: s.Explanation}, nil
	case StepFillBlank:
		return FillBlankPayload{Code: s.Code, Bank: s.Bank, Answers: s.Answers, Explanation: s.Explanation}, nil
	case StepArrange:
		return ArrangePayload{Prompt: s.Prompt, Items: s.Items, Order: s.Order, Explanation: s.Explanation}, nil
	case StepProse:
		if s.Rubric == nil {
			return nil, fmt.Errorf("step %q: prose requires a rubric", s.ID)
		}
		return ProsePayload{Prompt: s.Prompt, Rubric: *s.Rubric, MinChars: s.MinChars, Explanation: s.Explanation}, nil
	default:
		return nil, fmt.Errorf("step %q: unknown type %q", s.ID, s.Type)
	}
}

// ParsePayload decodes a stored payload into its typed form.
func ParsePayload(stepType string, raw []byte) (any, error) {
	switch stepType {
	case StepIntro:
		return decode[IntroPayload](stepType, raw)
	case StepMCQ:
		return decode[MCQPayload](stepType, raw)
	case StepTrueFalse:
		return decode[TrueFalsePayload](stepType, raw)
	case StepFillBlank:
		return decode[FillBlankPayload](stepType, raw)
	case StepArrange:
		return decode[ArrangePayload](stepType, raw)
	case StepProse:
		return decode[ProsePayload](stepType, raw)
	default:
		return nil, fmt.Errorf("unknown step type %q", stepType)
	}
}

func decode[T any](stepType string, raw []byte) (any, error) {
	var v T
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("decode %s payload: %w", stepType, err)
	}
	return v, nil
}

func asInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		if n == float64(int(n)) {
			return int(n), true
		}
	}
	return 0, false
}
