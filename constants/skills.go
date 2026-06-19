package constants

// Evaluation types supported by the question content model.
const (
	EvaluationTypeMultipleChoice = "multiple_choice"
	EvaluationTypeCoding           = "coding"
	EvaluationTypeSystemDesign     = "system_design"
	EvaluationTypeBehavioral       = "behavioral"
)

// DefaultEstimatedTimeMinutes is the fallback estimated completion time.
const DefaultEstimatedTimeMinutes = 10

// DefaultReadinessWeight is the baseline readiness contribution multiplier.
const DefaultReadinessWeight = 1.0
