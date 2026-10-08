-- Progress domain: weekly leagues. A learner joins the current week's league the first time
-- Progress grants them XP that week; XP earned that week accumulates on their membership.
-- Tiers are indexes into config.LeagueTiers (0 = lowest). Results are derived from final
-- standings when the learner next joins, so no job has to run at the week boundary.
CREATE TABLE league_cohorts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    week_start  DATE NOT NULL,
    tier        SMALLINT NOT NULL CHECK (tier >= 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Finding the oldest open cohort for a week and tier.
CREATE INDEX league_cohorts_week_tier_idx ON league_cohorts (week_start, tier, created_at);

CREATE TABLE league_memberships (
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    week_start  DATE NOT NULL,
    cohort_id   UUID NOT NULL REFERENCES league_cohorts(id) ON DELETE CASCADE,
    tier        SMALLINT NOT NULL CHECK (tier >= 0),
    weekly_xp   INT NOT NULL DEFAULT 0 CHECK (weekly_xp >= 0),
    last_xp_at  TIMESTAMPTZ NOT NULL,
    joined_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, week_start)
);

-- Standings: XP high to low, ties to whoever reached it first.
CREATE INDEX league_memberships_standings_idx ON league_memberships (cohort_id, weekly_xp DESC, last_xp_at);
