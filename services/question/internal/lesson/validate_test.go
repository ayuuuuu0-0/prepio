package lesson_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/prepio/prepio/services/question/internal/lesson"
	"github.com/stretchr/testify/require"
)

// validContent builds a minimal valid world + lesson. These are test fixtures only.
func validContent() *lesson.Content {
	return &lesson.Content{
		Worlds: []lesson.World{{
			Slug: "w-one", Name: "World One", Order: 1,
			Nodes: []lesson.Node{{Slug: "n-one", Label: "Node One", Type: lesson.NodeTypeLesson}},
		}},
		Lessons: []lesson.Lesson{{
			Slug: "l-one", Title: "Lesson One", Node: "n-one", Kind: lesson.KindLesson,
			Status: lesson.StatusPublished, Difficulty: "easy", EstMinutes: 5,
			Skills:  []lesson.SkillRef{{Skill: "s-one", Weight: 1.0}},
			Summary: []string{"one", "two"},
			Steps: []lesson.Step{
				{ID: "intro", Type: lesson.StepIntro, Beats: []lesson.Beat{{Text: "hello world", Emphasis: "world", DurationMs: 3000}}},
				{ID: "q1", Type: lesson.StepMCQ, Prompt: "p", Options: []string{"a", "b", "c"}, Answer: 1,
					Explanations: &lesson.Explanations{Correct: "b", Wrong: map[int]string{0: "no a", 2: "no c"}}},
				{ID: "q2", Type: lesson.StepTrueFalse, Statement: "s", Answer: true, Explanation: "e"},
				{ID: "q3", Type: lesson.StepFillBlank, Code: "x = ___1___", Bank: []string{"a", "b"}, Answers: []string{"a"}, Explanation: "e"},
			},
		}},
	}
}

var catalog = lesson.Catalog{Skills: map[string]bool{"s-one": true}}

func TestValidContentHasNoIssues(t *testing.T) {
	require.Empty(t, lesson.Validate(validContent(), catalog))
}

func TestValidationRules(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(c *lesson.Content)
		want   string
	}{
		{"unknown skill", func(c *lesson.Content) { c.Lessons[0].Skills[0].Skill = "ghost" }, `skill "ghost" does not exist`},
		{"weights must sum to 1", func(c *lesson.Content) { c.Lessons[0].Skills[0].Weight = 0.5 }, "must sum to 1.0"},
		{"no skills", func(c *lesson.Content) { c.Lessons[0].Skills = nil }, "at least one skill"},
		{"unknown node", func(c *lesson.Content) { c.Lessons[0].Node = "missing" }, `node "missing" does not exist`},
		{"node without lesson", func(c *lesson.Content) { c.Lessons = nil }, "has no lesson bound"},
		{"duplicate lesson slug", func(c *lesson.Content) { c.Lessons = append(c.Lessons, c.Lessons[0]) }, "duplicate lesson slug"},
		{"duplicate step id", func(c *lesson.Content) { c.Lessons[0].Steps[2].ID = "q1" }, "duplicate step id"},
		{"too few graded steps", func(c *lesson.Content) { c.Lessons[0].Steps = c.Lessons[0].Steps[:3] }, "needs 3-5 graded steps"},
		{"est minutes", func(c *lesson.Content) { c.Lessons[0].EstMinutes = 9 }, "est_minutes"},
		{"summary bullets", func(c *lesson.Content) { c.Lessons[0].Summary = []string{"only one"} }, "summary needs"},
		{"kind must match node type", func(c *lesson.Content) { c.Lessons[0].Kind = lesson.KindBoss }, "does not match node type"},
		{"intro must be first", func(c *lesson.Content) {
			s := c.Lessons[0].Steps
			s[0], s[1] = s[1], s[0]
		}, "intro must be the first step"},
		{"mcq needs 3-4 options", func(c *lesson.Content) { c.Lessons[0].Steps[1].Options = []string{"a", "b"} }, "3-4 options"},
		{"mcq answer in range", func(c *lesson.Content) { c.Lessons[0].Steps[1].Answer = 7 }, "answer must be an option index"},
		{"mcq wrong option needs why-not", func(c *lesson.Content) { delete(c.Lessons[0].Steps[1].Explanations.Wrong, 2) }, "wrong option 2"},
		{"mcq explains correct", func(c *lesson.Content) { c.Lessons[0].Steps[1].Explanations.Correct = "" }, "needs an explanation"},
		{"true_false explanation", func(c *lesson.Content) { c.Lessons[0].Steps[2].Explanation = "" }, "needs an explanation"},
		{"true_false answer is bool", func(c *lesson.Content) { c.Lessons[0].Steps[2].Answer = "yes" }, "must be true or false"},
		{"fill_blank blank count", func(c *lesson.Content) { c.Lessons[0].Steps[3].Answers = []string{"a", "b"} }, "1 blanks but 2 answers"},
		{"fill_blank bank has answers", func(c *lesson.Content) { c.Lessons[0].Steps[3].Bank = []string{"z"} }, `bank is missing answer "a"`},
		{"emphasis inside text", func(c *lesson.Content) { c.Lessons[0].Steps[0].Beats[0].Emphasis = "nope" }, "emphasis"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := validContent()
			tt.mutate(c)
			issues := lesson.Validate(c, catalog)
			require.NotEmpty(t, issues)
			require.Contains(t, strings.Join(issues, "\n"), tt.want)
		})
	}
}

func TestArrangeValidation(t *testing.T) {
	c := validContent()
	c.Lessons[0].Steps[3] = lesson.Step{ID: "q3", Type: lesson.StepArrange, Prompt: "p",
		Items: []string{"a", "b", "c"}, Order: []int{0, 1, 2}, Explanation: "e"}
	require.Contains(t, strings.Join(lesson.Validate(c, catalog), "\n"), "scramble")

	c.Lessons[0].Steps[3].Order = []int{2, 0, 0}
	require.Contains(t, strings.Join(lesson.Validate(c, catalog), "\n"), "permutation")

	c.Lessons[0].Steps[3].Order = []int{2, 0, 1}
	require.Empty(t, lesson.Validate(c, catalog))

	c.Lessons[0].Steps[3].Items = []string{"a", "b"}
	c.Lessons[0].Steps[3].Order = []int{1, 0}
	require.Contains(t, strings.Join(lesson.Validate(c, catalog), "\n"), "at least 3 items")
}

func TestProseValidation(t *testing.T) {
	c := validContent()
	c.Lessons[0].Steps[3] = lesson.Step{ID: "q3", Type: lesson.StepProse, Prompt: "p", MinChars: 150, Explanation: "e",
		Rubric: &lesson.Rubric{Concepts: []lesson.RubricConcept{{Name: "idea", Required: false}}}}
	require.Contains(t, strings.Join(lesson.Validate(c, catalog), "\n"), "at least one required concept")

	c.Lessons[0].Steps[3].Rubric.Concepts[0].Required = true
	require.Empty(t, lesson.Validate(c, catalog))

	c.Lessons[0].Steps[3].MinChars = 20
	require.Contains(t, strings.Join(lesson.Validate(c, catalog), "\n"), "minChars")
}

func TestNodeValidation(t *testing.T) {
	c := validContent()
	c.Worlds[0].Nodes = append(c.Worlds[0].Nodes, lesson.Node{Slug: "n-two", Label: "Two", Type: lesson.NodeTypeLesson, Requires: []string{"n-one"}})
	c.Worlds[0].Nodes[0].Requires = []string{"n-two"}
	issues := strings.Join(lesson.Validate(c, catalog), "\n")
	require.Contains(t, issues, "cycle")

	c = validContent()
	c.Worlds[0].Nodes[0].Requires = []string{"ghost"}
	require.Contains(t, strings.Join(lesson.Validate(c, catalog), "\n"), `requires unknown node "ghost"`)

	c = validContent()
	c.Worlds[0].Nodes[0].Slug = "Bad Slug"
	require.Contains(t, strings.Join(lesson.Validate(c, catalog), "\n"), "kebab-case")
}

func TestOfflineValidationSkipsSkillResolution(t *testing.T) {
	c := validContent()
	c.Lessons[0].Skills[0].Skill = "anything"
	require.Empty(t, lesson.Validate(c, lesson.Catalog{}))
}

func TestHashIgnoresStatusButNotContent(t *testing.T) {
	a := validContent().Lessons[0]
	b := validContent().Lessons[0]
	b.Status = lesson.StatusDraft
	ha, err := a.Hash()
	require.NoError(t, err)
	hb, err := b.Hash()
	require.NoError(t, err)
	require.Equal(t, ha, hb)

	b.Steps[2].Explanation = "changed"
	hc, err := b.Hash()
	require.NoError(t, err)
	require.NotEqual(t, ha, hc)
}

func TestProjectNeverLeaksAnswers(t *testing.T) {
	const secret = "SECRET-SENTINEL"
	steps := []struct {
		typ     string
		payload any
	}{
		{lesson.StepMCQ, lesson.MCQPayload{Prompt: "p", Options: []string{"a", "b", "c"}, Answer: 2,
			Explanations: lesson.Explanations{Correct: secret, Wrong: map[int]string{0: secret, 1: secret}}}},
		{lesson.StepTrueFalse, lesson.TrueFalsePayload{Statement: "s", Answer: true, Explanation: secret}},
		{lesson.StepFillBlank, lesson.FillBlankPayload{Code: "___1___", Bank: []string{"x", "y"}, Answers: []string{"x"}, Explanation: secret}},
		{lesson.StepArrange, lesson.ArrangePayload{Prompt: "p", Items: []string{"a", "b", "c"}, Order: []int{2, 0, 1}, Explanation: secret}},
		{lesson.StepProse, lesson.ProsePayload{Prompt: "p", MinChars: 100, Explanation: secret,
			Rubric: lesson.Rubric{Concepts: []lesson.RubricConcept{{Name: secret, Aliases: []string{secret}, Required: true}}}}},
	}
	for _, s := range steps {
		t.Run(s.typ, func(t *testing.T) {
			projected, err := lesson.Project("k", s.typ, 0, s.payload, "seed")
			require.NoError(t, err)
			raw, err := json.Marshal(projected)
			require.NoError(t, err)
			body := string(raw)
			require.NotContains(t, body, secret)
			for _, field := range []string{`"answer"`, `"answers"`, `"order"`, `"explanation`, `"rubric"`, `"correct"`, `"wrong"`} {
				require.NotContains(t, body, field)
			}
		})
	}
}

func TestFillBlankBankOrderIsStablePerSeed(t *testing.T) {
	p := lesson.FillBlankPayload{Code: "___1___", Bank: []string{"a", "b", "c", "d", "e", "f"}, Answers: []string{"a"}}
	first, err := lesson.Project("k", lesson.StepFillBlank, 0, p, "attempt-1")
	require.NoError(t, err)
	again, err := lesson.Project("k", lesson.StepFillBlank, 0, p, "attempt-1")
	require.NoError(t, err)
	require.Equal(t, first.FillBlank.Bank, again.FillBlank.Bank)
	require.ElementsMatch(t, p.Bank, first.FillBlank.Bank)
}
