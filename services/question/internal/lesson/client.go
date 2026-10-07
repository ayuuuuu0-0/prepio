package lesson

import (
	"fmt"
	"hash/fnv"
	"math/rand"
)

// ClientStep is the only step shape ever sent to a client. It is built from an
// explicit whitelist per step type, so answers, explanations, and rubrics can
// never leak by adding a field to a stored payload.
type ClientStep struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Position int    `json:"position"`

	Intro     *ClientIntro     `json:"intro,omitempty"`
	MCQ       *ClientMCQ       `json:"mcq,omitempty"`
	TrueFalse *ClientTrueFalse `json:"true_false,omitempty"`
	FillBlank *ClientFillBlank `json:"fill_blank,omitempty"`
	Arrange   *ClientArrange   `json:"arrange,omitempty"`
	Prose     *ClientProse     `json:"prose,omitempty"`
}

// ClientIntro is the client view of an intro step.
type ClientIntro struct {
	Beats    []Beat `json:"beats"`
	MediaURL string `json:"media_url,omitempty"`
}

// ClientMCQ is the client view of a multiple-choice step.
type ClientMCQ struct {
	Prompt  string   `json:"prompt"`
	Options []string `json:"options"`
}

// ClientTrueFalse is the client view of a true/false step.
type ClientTrueFalse struct {
	Statement string `json:"statement"`
}

// ClientFillBlank is the client view of a fill-in-the-blank step.
type ClientFillBlank struct {
	Code   string   `json:"code"`
	Blanks int      `json:"blanks"`
	Bank   []string `json:"bank"`
}

// ClientArrange is the client view of an arrange-in-order step.
type ClientArrange struct {
	Prompt string   `json:"prompt"`
	Items  []string `json:"items"`
}

// ClientProse is the client view of a prose step.
type ClientProse struct {
	Prompt   string `json:"prompt"`
	MinChars int    `json:"min_chars"`
}

// Project builds the answer-free client view of a stored step payload.
// seed makes the fill-blank bank order stable for one attempt.
func Project(stepKey, stepType string, position int, payload any, seed string) (ClientStep, error) {
	out := ClientStep{ID: stepKey, Type: stepType, Position: position}
	switch p := payload.(type) {
	case IntroPayload:
		out.Intro = &ClientIntro{Beats: p.Beats, MediaURL: p.MediaURL}
	case MCQPayload:
		out.MCQ = &ClientMCQ{Prompt: p.Prompt, Options: p.Options}
	case TrueFalsePayload:
		out.TrueFalse = &ClientTrueFalse{Statement: p.Statement}
	case FillBlankPayload:
		out.FillBlank = &ClientFillBlank{Code: p.Code, Blanks: len(p.Answers), Bank: shuffled(p.Bank, seed+":"+stepKey)}
	case ArrangePayload:
		out.Arrange = &ClientArrange{Prompt: p.Prompt, Items: p.Items}
	case ProsePayload:
		out.Prose = &ClientProse{Prompt: p.Prompt, MinChars: p.MinChars}
	default:
		return ClientStep{}, fmt.Errorf("step %q: unsupported payload %T", stepKey, payload)
	}
	return out, nil
}

func shuffled(in []string, seed string) []string {
	out := append([]string(nil), in...)
	h := fnv.New64a()
	_, _ = h.Write([]byte(seed))
	r := rand.New(rand.NewSource(int64(h.Sum64())))
	r.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })
	return out
}
