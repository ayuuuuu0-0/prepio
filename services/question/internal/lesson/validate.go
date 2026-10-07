package lesson

import (
	"fmt"
	"math"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/prepio/prepio/constants"
)

var (
	slugPattern  = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	blankPattern = regexp.MustCompile(`___(\d+)___`)
)

const (
	maxLessonMinutes = 5
	maxBossMinutes   = 8
	minGradedSteps   = 3
	maxGradedLesson  = 5
	maxGradedBoss    = 8
	minSummary       = 2
	maxSummary       = 4
	weightTolerance  = 0.001

	// Intro timing: a beat is long enough to read and short enough not to drag, and the whole
	// intro stays quick because it is always skippable but should not need to be.
	minBeatMs  = 1500
	maxBeatMs  = 12000
	maxIntroMs = 45000
)

// IntroVisuals are the visual names an intro beat may use. Each must exist in the web client
// (web/src/components/lesson/visuals.tsx); the list is kept here so an unknown name fails CI
// instead of silently rendering nothing.
var IntroVisuals = []string{"bars", "flow", "pulse"}

var mediaExtensions = []string{".mp4", ".webm", ".png", ".jpg", ".jpeg", ".gif", ".webp", ".avif", ".svg"}

// validMediaURL accepts https URLs and site-relative paths that point at known media types.
// It mirrors the client, which refuses to render anything else.
func validMediaURL(raw string) bool {
	if strings.HasPrefix(raw, "//") {
		return false
	}
	path := raw
	if strings.HasPrefix(raw, "/") {
		path = raw
	} else {
		u, err := url.Parse(raw)
		if err != nil || u.Scheme != "https" || u.Host == "" {
			return false
		}
		path = u.Path
	}
	if i := strings.IndexAny(path, "?#"); i >= 0 {
		path = path[:i]
	}
	path = strings.ToLower(path)
	for _, ext := range mediaExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

// Catalog is the set of references a lesson may point to.
// A nil Skills map skips skill resolution (offline validation).
type Catalog struct {
	Skills map[string]bool
}

// Validate checks all worlds and lessons and returns every problem found.
// An empty result means the content can be published.
func Validate(c *Content, catalog Catalog) []string {
	v := &validator{catalog: catalog}
	v.worlds(c)
	v.lessons(c)
	v.binding(c)
	return v.issues
}

type validator struct {
	catalog Catalog
	issues  []string
}

func (v *validator) addf(scope, format string, args ...any) {
	v.issues = append(v.issues, scope+": "+fmt.Sprintf(format, args...))
}

func (v *validator) worlds(c *Content) {
	worldSlugs := map[string]bool{}
	nodeSlugs := map[string]bool{}
	for _, w := range c.Worlds {
		scope := "world " + w.Slug
		if !slugPattern.MatchString(w.Slug) {
			v.addf(scope, "slug must be lowercase kebab-case")
		}
		if worldSlugs[w.Slug] {
			v.addf(scope, "duplicate world slug")
		}
		worldSlugs[w.Slug] = true
		if strings.TrimSpace(w.Name) == "" {
			v.addf(scope, "name is required")
		}
		if len(w.Nodes) == 0 {
			v.addf(scope, "a world needs at least one node")
		}
		for _, n := range w.Nodes {
			nscope := scope + " node " + n.Slug
			if !slugPattern.MatchString(n.Slug) {
				v.addf(nscope, "slug must be lowercase kebab-case")
			}
			if nodeSlugs[n.Slug] {
				v.addf(nscope, "duplicate node slug (node slugs are unique across worlds)")
			}
			nodeSlugs[n.Slug] = true
			if strings.TrimSpace(n.Label) == "" {
				v.addf(nscope, "label is required")
			}
			if n.Type != NodeTypeLesson && n.Type != NodeTypeBoss {
				v.addf(nscope, "type must be %q or %q", NodeTypeLesson, NodeTypeBoss)
			}
		}
	}

	requires := map[string][]string{}
	for _, w := range c.Worlds {
		for _, n := range w.Nodes {
			for _, r := range n.Requires {
				switch {
				case r == n.Slug:
					v.addf("node "+n.Slug, "cannot require itself")
				case !nodeSlugs[r]:
					v.addf("node "+n.Slug, "requires unknown node %q", r)
				default:
					requires[n.Slug] = append(requires[n.Slug], r)
				}
			}
		}
	}
	for slug := range nodeSlugs {
		if hasCycle(slug, requires, map[string]int{}) {
			v.addf("node "+slug, "prerequisites form a cycle")
		}
	}
}

func hasCycle(node string, edges map[string][]string, state map[string]int) bool {
	switch state[node] {
	case 1:
		return true
	case 2:
		return false
	}
	state[node] = 1
	for _, next := range edges[node] {
		if hasCycle(next, edges, state) {
			return true
		}
	}
	state[node] = 2
	return false
}

// binding checks that every node binds exactly one lesson and every lesson binds a real node.
func (v *validator) binding(c *Content) {
	nodeTypes := map[string]string{}
	for _, w := range c.Worlds {
		for _, n := range w.Nodes {
			nodeTypes[n.Slug] = n.Type
		}
	}
	bound := map[string]string{}
	for _, l := range c.Lessons {
		scope := "lesson " + l.Slug
		nodeType, ok := nodeTypes[l.Node]
		if !ok {
			v.addf(scope, "node %q does not exist", l.Node)
			continue
		}
		if other, dup := bound[l.Node]; dup {
			v.addf(scope, "node %q is already bound by lesson %q", l.Node, other)
		}
		bound[l.Node] = l.Slug
		if nodeType != l.Kind {
			v.addf(scope, "kind %q does not match node type %q", l.Kind, nodeType)
		}
	}
	for node := range nodeTypes {
		if _, ok := bound[node]; !ok {
			v.addf("node "+node, "has no lesson bound to it")
		}
	}
}

func (v *validator) lessons(c *Content) {
	slugs := map[string]bool{}
	for i := range c.Lessons {
		l := &c.Lessons[i]
		scope := "lesson " + l.Slug
		if !slugPattern.MatchString(l.Slug) {
			v.addf(scope, "slug must be lowercase kebab-case")
		}
		if slugs[l.Slug] {
			v.addf(scope, "duplicate lesson slug")
		}
		slugs[l.Slug] = true
		v.lesson(scope, l)
	}
}

func (v *validator) lesson(scope string, l *Lesson) {
	if strings.TrimSpace(l.Title) == "" {
		v.addf(scope, "title is required")
	}
	if !slices.Contains([]string{KindLesson, KindBoss}, l.Kind) {
		v.addf(scope, "kind must be %q or %q", KindLesson, KindBoss)
	}
	if !slices.Contains([]string{StatusDraft, StatusReview, StatusPublished, StatusDeprecated}, l.Status) {
		v.addf(scope, "status must be draft, review, published, or deprecated")
	}
	if !slices.Contains([]string{"easy", "medium", "hard"}, l.Difficulty) {
		v.addf(scope, "difficulty must be easy, medium, or hard")
	}
	limit := maxLessonMinutes
	if l.Kind == KindBoss {
		limit = maxBossMinutes
	}
	if l.EstMinutes < 1 || l.EstMinutes > limit {
		v.addf(scope, "est_minutes must be between 1 and %d for kind %q", limit, l.Kind)
	}
	if len(l.Summary) < minSummary || len(l.Summary) > maxSummary {
		v.addf(scope, "summary needs %d-%d takeaway bullets", minSummary, maxSummary)
	}
	for _, s := range l.Summary {
		if strings.TrimSpace(s) == "" {
			v.addf(scope, "summary bullets cannot be empty")
		}
	}

	v.skills(scope, l)
	v.steps(scope, l)
}

func (v *validator) skills(scope string, l *Lesson) {
	if len(l.Skills) == 0 {
		v.addf(scope, "every lesson maps to at least one skill")
		return
	}
	seen := map[string]bool{}
	sum := 0.0
	for _, s := range l.Skills {
		if seen[s.Skill] {
			v.addf(scope, "skill %q listed twice", s.Skill)
		}
		seen[s.Skill] = true
		if s.Weight <= 0 || s.Weight > 1 {
			v.addf(scope, "skill %q weight must be in (0, 1]", s.Skill)
		}
		if v.catalog.Skills != nil && !v.catalog.Skills[s.Skill] {
			v.addf(scope, "skill %q does not exist in the catalog", s.Skill)
		}
		sum += s.Weight
	}
	if math.Abs(sum-1.0) > weightTolerance {
		v.addf(scope, "skill weights sum to %.3f, must sum to 1.0", sum)
	}
}

func (v *validator) steps(scope string, l *Lesson) {
	ids := map[string]bool{}
	graded := 0
	for i, s := range l.Steps {
		sscope := scope + " step " + s.ID
		if !slugPattern.MatchString(s.ID) {
			v.addf(sscope, "step id must be lowercase kebab-case")
		}
		if ids[s.ID] {
			v.addf(sscope, "duplicate step id")
		}
		ids[s.ID] = true

		if s.Type == StepIntro {
			if i != 0 {
				v.addf(sscope, "an intro must be the first step, and a lesson has at most one")
			}
		} else {
			graded++
		}
		v.step(sscope, s)
	}

	maxGraded := maxGradedLesson
	if l.Kind == KindBoss {
		maxGraded = maxGradedBoss
	}
	if graded < minGradedSteps || graded > maxGraded {
		v.addf(scope, "needs %d-%d graded steps, has %d", minGradedSteps, maxGraded, graded)
	}
}

func (v *validator) step(scope string, s Step) {
	if s.Type != StepIntro && strings.TrimSpace(explanationOf(s)) == "" && s.Type != StepMCQ {
		v.addf(scope, "every graded step needs an explanation")
	}
	switch s.Type {
	case StepIntro:
		v.intro(scope, s)
	case StepMCQ:
		v.mcq(scope, s)
	case StepTrueFalse:
		v.trueFalse(scope, s)
	case StepFillBlank:
		v.fillBlank(scope, s)
	case StepArrange:
		v.arrange(scope, s)
	case StepProse:
		v.prose(scope, s)
	default:
		v.addf(scope, "unknown step type %q", s.Type)
	}
}

func explanationOf(s Step) string {
	if s.Type == StepMCQ && s.Explanations != nil {
		return s.Explanations.Correct
	}
	return s.Explanation
}

func (v *validator) intro(scope string, s Step) {
	if len(s.Beats) == 0 {
		v.addf(scope, "intro needs at least one beat")
	}
	if s.MediaURL != "" && !validMediaURL(s.MediaURL) {
		v.addf(scope, "mediaUrl must be an https URL or a site-relative path to a known image or video type")
	}
	total := 0
	for i, b := range s.Beats {
		total += b.DurationMs
		if strings.TrimSpace(b.Text) == "" {
			v.addf(scope, "beat %d has no text", i+1)
		}
		if b.DurationMs < minBeatMs || b.DurationMs > maxBeatMs {
			v.addf(scope, "beat %d durationMs must be between %d and %d", i+1, minBeatMs, maxBeatMs)
		}
		if b.Visual != "" && !slices.Contains(IntroVisuals, b.Visual) {
			v.addf(scope, "beat %d visual %q is unknown (allowed: %s)", i+1, b.Visual, strings.Join(IntroVisuals, ", "))
		}
		if b.Emphasis != "" && !strings.Contains(b.Text, b.Emphasis) {
			v.addf(scope, "beat %d emphasis %q is not part of its text", i+1, b.Emphasis)
		}
	}
	if total > maxIntroMs {
		v.addf(scope, "intro runs %ds, keep it under %ds", total/1000, maxIntroMs/1000)
	}
}

func (v *validator) mcq(scope string, s Step) {
	if strings.TrimSpace(s.Prompt) == "" {
		v.addf(scope, "mcq needs a prompt")
	}
	if len(s.Options) < 3 || len(s.Options) > 4 {
		v.addf(scope, "mcq needs 3-4 options, has %d", len(s.Options))
	}
	seen := map[string]bool{}
	for _, o := range s.Options {
		if strings.TrimSpace(o) == "" {
			v.addf(scope, "mcq options cannot be empty")
		}
		if seen[o] {
			v.addf(scope, "mcq option %q is duplicated", o)
		}
		seen[o] = true
	}
	answer, ok := asInt(s.Answer)
	if !ok || answer < 0 || answer >= len(s.Options) {
		v.addf(scope, "mcq answer must be an option index (0-%d)", len(s.Options)-1)
		return
	}
	if s.Explanations == nil || strings.TrimSpace(s.Explanations.Correct) == "" {
		v.addf(scope, "every graded step needs an explanation (explanations.correct)")
		return
	}
	for i := range s.Options {
		if i == answer {
			if _, extra := s.Explanations.Wrong[i]; extra {
				v.addf(scope, "explanations.wrong must not include the correct option %d", i)
			}
			continue
		}
		if strings.TrimSpace(s.Explanations.Wrong[i]) == "" {
			v.addf(scope, "mcq needs a 'why not' explanation for wrong option %d", i)
		}
	}
	for i := range s.Explanations.Wrong {
		if i < 0 || i >= len(s.Options) {
			v.addf(scope, "explanations.wrong has index %d outside the options", i)
		}
	}
}

func (v *validator) trueFalse(scope string, s Step) {
	if strings.TrimSpace(s.Statement) == "" {
		v.addf(scope, "true_false needs a statement")
	}
	if _, ok := s.Answer.(bool); !ok {
		v.addf(scope, "true_false answer must be true or false")
	}
}

func (v *validator) fillBlank(scope string, s Step) {
	if strings.TrimSpace(s.Code) == "" {
		v.addf(scope, "fill_blank needs code")
		return
	}
	found := map[int]int{}
	for _, m := range blankPattern.FindAllStringSubmatch(s.Code, -1) {
		n, _ := strconv.Atoi(m[1])
		found[n]++
	}
	if len(found) == 0 {
		v.addf(scope, "fill_blank code needs numbered blanks like ___1___")
	}
	for n, count := range found {
		if count > 1 {
			v.addf(scope, "blank %d appears more than once", n)
		}
		if n < 1 || n > len(found) {
			v.addf(scope, "blanks must be numbered 1..%d, found %d", len(found), n)
		}
	}
	if len(found) != len(s.Answers) {
		v.addf(scope, "code has %d blanks but %d answers", len(found), len(s.Answers))
	}
	remaining := map[string]int{}
	for _, b := range s.Bank {
		if strings.TrimSpace(b) == "" {
			v.addf(scope, "bank entries cannot be empty")
		}
		remaining[b]++
	}
	for _, a := range s.Answers {
		if strings.TrimSpace(a) == "" {
			v.addf(scope, "answers cannot be empty")
		}
		if remaining[a] == 0 {
			v.addf(scope, "bank is missing answer %q", a)
			continue
		}
		remaining[a]--
	}
}

func (v *validator) arrange(scope string, s Step) {
	if strings.TrimSpace(s.Prompt) == "" {
		v.addf(scope, "arrange needs a prompt")
	}
	if len(s.Items) < 3 {
		v.addf(scope, "arrange needs at least 3 items, has %d", len(s.Items))
	}
	for _, it := range s.Items {
		if strings.TrimSpace(it) == "" {
			v.addf(scope, "arrange items cannot be empty")
		}
	}
	if len(s.Order) != len(s.Items) {
		v.addf(scope, "order must list every item index once")
		return
	}
	seen := make([]bool, len(s.Items))
	identity := true
	for pos, idx := range s.Order {
		if idx < 0 || idx >= len(s.Items) || seen[idx] {
			v.addf(scope, "order must be a permutation of the item indices")
			return
		}
		seen[idx] = true
		if idx != pos {
			identity = false
		}
	}
	if identity {
		v.addf(scope, "items are authored in the correct order; scramble them so the client does not receive the answer")
	}
}

func (v *validator) prose(scope string, s Step) {
	if strings.TrimSpace(s.Prompt) == "" {
		v.addf(scope, "prose needs a prompt")
	}
	if s.MinChars < constants.MinAnswerLength {
		v.addf(scope, "prose minChars must be at least %d", constants.MinAnswerLength)
	}
	if s.Rubric == nil {
		v.addf(scope, "prose needs a rubric")
		return
	}
	required := 0
	for _, c := range s.Rubric.Concepts {
		if strings.TrimSpace(c.Name) == "" {
			v.addf(scope, "rubric concepts need a name")
		}
		if c.Required {
			required++
		}
	}
	if required == 0 {
		v.addf(scope, "rubric needs at least one required concept")
	}
}
