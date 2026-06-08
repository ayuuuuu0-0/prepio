package config

// LevelThresholds holds cumulative XP required to reach each level.
// Index 0 is level 1 (0 XP), index 1 is level 2, and so on.
var LevelThresholds = []int{
	0,    // level 1
	100,  // level 2
	250,  // level 3
	500,  // level 4
	800,  // level 5
	1200, // level 6
	1700, // level 7
	2300, // level 8
	3000, // level 9
	3800, // level 10
}

// CurrentLevel returns the level for a given total XP amount.
func CurrentLevel(totalXP int) int {
	level := 1
	for i := len(LevelThresholds) - 1; i >= 0; i-- {
		if totalXP >= LevelThresholds[i] {
			level = i + 1
			break
		}
	}
	return level
}

// DifficultyForLevel maps a user level to a question difficulty band.
func DifficultyForLevel(level int) string {
	if level <= 2 {
		return "easy"
	}
	if level <= 5 {
		return "medium"
	}
	return "hard"
}

// XPToNextLevel returns XP remaining until the next level, or 0 at max level.
func XPToNextLevel(totalXP int) int {
	level := CurrentLevel(totalXP)
	if level >= len(LevelThresholds) {
		return 0
	}
	return LevelThresholds[level] - totalXP
}
