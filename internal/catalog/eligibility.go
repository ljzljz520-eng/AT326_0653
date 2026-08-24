package catalog

import (
	"fmt"
	"strings"

	"arcadeledger/internal/model"
)

func NormalizeTier(value string) string { return strings.ToLower(strings.TrimSpace(value)) }

func ValidTier(value string) bool {
	switch NormalizeTier(value) {
	case "bronze", "silver", "gold", "champion":
		return true
	default:
		return false
	}
}

func ValidateClaim(prize model.PrizeDescription) error {
	if err := model.ValidatePrize(prize); err != nil {
		return err
	}
	if !prize.Active {
		return fmt.Errorf("prize is inactive")
	}
	if !ValidTier(prize.Tier) {
		return fmt.Errorf("unknown prize tier")
	}
	if prize.Claimed >= prize.Quantity {
		return model.ErrPrizeExhausted
	}
	return nil
}

func AwardRank(rank int) string {
	if rank <= 1 {
		return "champion"
	}
	if rank <= 3 {
		return "gold"
	}
	if rank <= 10 {
		return "silver"
	}
	return "bronze"
}
