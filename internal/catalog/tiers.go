package catalog

import (
	"sort"
	"strings"

	"arcadeledger/internal/model"
)

var tierOrder = map[string]int{"bronze": 1, "silver": 2, "gold": 3, "champion": 4}

func TierScore(tier string) int { return tierOrder[NormalizeTier(tier)] }

func CompareTiers(left, right string) int {
	l, r := TierScore(left), TierScore(right)
	if l < r {
		return -1
	}
	if l > r {
		return 1
	}
	return 0
}

func SortPrizes(prizes []model.PrizeDescription) []model.PrizeDescription {
	result := append([]model.PrizeDescription(nil), prizes...)
	sort.SliceStable(result, func(i, j int) bool {
		if CompareTiers(result[i].Tier, result[j].Tier) == 0 {
			return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
		}
		return CompareTiers(result[i].Tier, result[j].Tier) > 0
	})
	return result
}

func Available(prizes []model.PrizeDescription, tier string) []model.PrizeDescription {
	result := make([]model.PrizeDescription, 0)
	for _, prize := range prizes {
		if prize.Active && strings.EqualFold(prize.Tier, tier) && prize.Claimed < prize.Quantity {
			result = append(result, prize)
		}
	}
	return SortPrizes(result)
}

func SelectBest(prizes []model.PrizeDescription, rank int) (model.PrizeDescription, bool) {
	tier := AwardRank(rank)
	available := Available(prizes, tier)
	if len(available) == 0 {
		return model.PrizeDescription{}, false
	}
	return available[0], true
}

func TierLabels() []string { return []string{"bronze", "silver", "gold", "champion"} }
