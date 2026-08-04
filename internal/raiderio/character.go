package raiderio

type Character struct {
	Name           string `json:"name"`
	Realm          string `json:"realm"`
	Class          string `json:"class"`
	ActiveSpecName string `json:"active_spec_name"`

	Gear struct {
		ItemLevelEquipped float64 `json:"item_level_equipped"`
		ItemLevelTotal    float64 `json:"item_level_total"`
	} `json:"gear"`

	MythicPlusScores []struct {
		Scores struct {
			All float64 `json:"all"`
		} `json:"scores"`
	} `json:"mythic_plus_scores_by_season"`

	RaidProgression map[string]RaidProgression `json:"raid_progression"`
}

type RaidProgression struct {
	Summary            string `json:"summary"`
	TotalBosses        int    `json:"total_bosses"`
	NormalBossesKilled int    `json:"normal_bosses_killed"`
	HeroicBossesKilled int    `json:"heroic_bosses_killed"`
	MythicBossesKilled int    `json:"mythic_bosses_killed"`
}

func (c *Character) ItemLevel() float64 {
	return c.Gear.ItemLevelEquipped
}

func (c *Character) MythicRating() float64 {
	if len(c.MythicPlusScores) == 0 {
		return 0
	}

	return c.MythicPlusScores[0].Scores.All
}
