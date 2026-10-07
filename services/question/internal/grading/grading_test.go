package grading_test

import (
	"errors"
	"testing"

	"github.com/prepio/prepio/services/question/internal/grading"
	"github.com/prepio/prepio/services/question/internal/lesson"
	"github.com/stretchr/testify/require"
)

func intp(v int) *int    { return &v }
func boolp(v bool) *bool { return &v }

var mcq = lesson.MCQPayload{
	Prompt:  "p",
	Options: []string{"a", "b", "c"},
	Answer:  1,
	Explanations: lesson.Explanations{
		Correct: "b works",
		Wrong:   map[int]string{0: "not a", 2: "not c"},
	},
}

func TestMCQ(t *testing.T) {
	tests := []struct {
		name       string
		answer     grading.Answer
		wantErr    bool
		wantOK     bool
		wantWhyNot string
	}{
		{name: "correct", answer: grading.Answer{Choice: intp(1)}, wantOK: true},
		{name: "wrong first option", answer: grading.Answer{Choice: intp(0)}, wantWhyNot: "not a"},
		{name: "wrong last option", answer: grading.Answer{Choice: intp(2)}, wantWhyNot: "not c"},
		{name: "missing choice", answer: grading.Answer{}, wantErr: true},
		{name: "negative index", answer: grading.Answer{Choice: intp(-1)}, wantErr: true},
		{name: "index out of range", answer: grading.Answer{Choice: intp(3)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := grading.Grade(mcq, tt.answer)
			if tt.wantErr {
				require.ErrorIs(t, err, grading.ErrInvalidAnswer)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantOK, res.Correct)
			require.Equal(t, "b works", res.Explanation)
			require.Equal(t, tt.wantWhyNot, res.WhyNot)
			if tt.wantOK {
				require.Nil(t, res.CorrectAnswer, "correct answer is only revealed after an incorrect answer")
			} else {
				require.NotNil(t, res.CorrectAnswer)
				require.Equal(t, 1, *res.CorrectAnswer.Choice)
			}
		})
	}
}

func TestTrueFalse(t *testing.T) {
	p := lesson.TrueFalsePayload{Statement: "s", Answer: false, Explanation: "because"}

	res, err := grading.Grade(p, grading.Answer{Value: boolp(false)})
	require.NoError(t, err)
	require.True(t, res.Correct)
	require.Nil(t, res.CorrectAnswer)

	res, err = grading.Grade(p, grading.Answer{Value: boolp(true)})
	require.NoError(t, err)
	require.False(t, res.Correct)
	require.False(t, *res.CorrectAnswer.Value)
	require.Equal(t, "because", res.Explanation)

	_, err = grading.Grade(p, grading.Answer{})
	require.ErrorIs(t, err, grading.ErrInvalidAnswer)
}

func TestFillBlank(t *testing.T) {
	p := lesson.FillBlankPayload{
		Code:        "x = ___1___ + ___2___",
		Bank:        []string{"a", "b", "z"},
		Answers:     []string{"a", "b"},
		Explanation: "e",
	}
	tests := []struct {
		name    string
		blanks  []string
		wantOK  bool
		wantErr bool
	}{
		{name: "correct", blanks: []string{"a", "b"}, wantOK: true},
		{name: "whitespace is trimmed", blanks: []string{" a ", "b\n"}, wantOK: true},
		{name: "swapped", blanks: []string{"b", "a"}},
		{name: "distractor", blanks: []string{"a", "z"}},
		{name: "case matters", blanks: []string{"A", "b"}},
		{name: "too few blanks", blanks: []string{"a"}, wantErr: true},
		{name: "too many blanks", blanks: []string{"a", "b", "c"}, wantErr: true},
		{name: "no blanks", blanks: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := grading.Grade(p, grading.Answer{Blanks: tt.blanks})
			if tt.wantErr {
				require.ErrorIs(t, err, grading.ErrInvalidAnswer)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantOK, res.Correct)
			if tt.wantOK {
				require.Nil(t, res.CorrectAnswer)
			} else {
				require.Equal(t, []string{"a", "b"}, res.CorrectAnswer.Blanks)
			}
		})
	}
}

func TestArrange(t *testing.T) {
	p := lesson.ArrangePayload{Prompt: "p", Items: []string{"c", "a", "b"}, Order: []int{1, 2, 0}, Explanation: "e"}
	tests := []struct {
		name    string
		order   []int
		wantOK  bool
		wantErr bool
	}{
		{name: "correct", order: []int{1, 2, 0}, wantOK: true},
		{name: "wrong permutation", order: []int{0, 1, 2}},
		{name: "duplicate index", order: []int{1, 1, 0}, wantErr: true},
		{name: "out of range", order: []int{1, 2, 3}, wantErr: true},
		{name: "negative", order: []int{-1, 2, 0}, wantErr: true},
		{name: "wrong length", order: []int{1, 2}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := grading.Grade(p, grading.Answer{Order: tt.order})
			if tt.wantErr {
				require.ErrorIs(t, err, grading.ErrInvalidAnswer)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantOK, res.Correct)
			if !tt.wantOK {
				require.Equal(t, []int{1, 2, 0}, res.CorrectAnswer.Order)
			}
		})
	}
}

func TestProseIsNotGradedHere(t *testing.T) {
	_, err := grading.Grade(lesson.ProsePayload{}, grading.Answer{})
	require.True(t, errors.Is(err, grading.ErrNotDeterministic))
	_, err = grading.Grade(lesson.IntroPayload{}, grading.Answer{})
	require.True(t, errors.Is(err, grading.ErrNotDeterministic))
}
