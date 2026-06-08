package constants

// Starter companion IDs available during onboarding (vision.md).
const (
	CompanionByteID = "b0000000-0000-4000-8000-000000000001"
	CompanionPipID  = "b0000000-0000-4000-8000-000000000002"
	CompanionNovaID = "b0000000-0000-4000-8000-000000000003"
	CompanionKodoID = "b0000000-0000-4000-8000-000000000004"
	CompanionZaraID = "b0000000-0000-4000-8000-000000000005"
)

// DefaultCharacterID is the fallback companion before onboarding completes.
const DefaultCharacterID = CompanionByteID

// StarterCompanionIDs lists companions selectable during onboarding.
var StarterCompanionIDs = []string{
	CompanionByteID,
	CompanionPipID,
	CompanionNovaID,
	CompanionKodoID,
	CompanionZaraID,
}
