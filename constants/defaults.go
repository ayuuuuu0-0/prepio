package constants

// DefaultTimezone is the default user timezone.
const DefaultTimezone = "Asia/Kolkata"

// DefaultReminderTime is the default streak reminder time (HH:MM:SS).
const DefaultReminderTime = "21:50:00"

// MinAnswerLength is the minimum trimmed characters required before submitting an answer.
const MinAnswerLength = 100

// MinEvaluationScore is the minimum score (0–100) for an answer to count as correct.
const MinEvaluationScore = 60

// MinKeywordFilterRatio is the minimum ratio of matched required concepts (0–1)
// for an answer to pass Stage 1 and proceed to Stage 2 structural evaluation.
// Answers below this threshold are rejected immediately without an LLM call.
const MinKeywordFilterRatio = 0.35
