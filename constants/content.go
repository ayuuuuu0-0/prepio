package constants

// Pool selection strategies for node_pools.selection_strategy.
const (
	PoolSelectionRandomUnseen = "random_unseen"
	PoolSelectionSequential   = "sequential"
	PoolSelectionBossMixed    = "boss_mixed"
)

// EnvJourneyPoolSelection enables pool-based question selection in journey (Journey V2).
const EnvJourneyPoolSelection = "JOURNEY_POOL_SELECTION"

// EnvReadinessV2 enables V2 readiness visibility on the dashboard.
const EnvReadinessV2 = "READINESS_V2"

// FoundationForestWorldSlug is the first journey world.
const FoundationForestWorldSlug = "foundation-forest"
