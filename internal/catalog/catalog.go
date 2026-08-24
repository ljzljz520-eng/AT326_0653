package catalog

import (
	"fmt"
	"strings"

	"arcadeledger/internal/model"
	"arcadeledger/internal/store"
)

type Catalog struct{ db *store.DB }

func New(db *store.DB) *Catalog { return &Catalog{db: db} }

func (c *Catalog) Add(id, name, tier, description string, quantity int) (model.PrizeDescription, error) {
	if c == nil || c.db == nil {
		return model.PrizeDescription{}, fmt.Errorf("catalog unavailable")
	}
	prize := model.PrizeDescription{ID: strings.TrimSpace(id), Name: strings.TrimSpace(name), Tier: NormalizeTier(tier), Description: model.NormalizeDescription(description), Quantity: quantity, Active: true}
	if err := model.ValidatePrize(prize); err != nil {
		return model.PrizeDescription{}, err
	}
	if !ValidTier(prize.Tier) {
		return model.PrizeDescription{}, fmt.Errorf("unknown prize tier")
	}
	if _, err := c.db.GetPrize(prize.ID); err == nil {
		return model.PrizeDescription{}, fmt.Errorf("prize already exists")
	}
	return prize, c.db.SavePrize(prize)
}

func (c *Catalog) Get(id string) (model.PrizeDescription, error) { return c.db.GetPrize(id) }

func (c *Catalog) List(activeOnly bool) ([]model.PrizeDescription, error) {
	return c.db.ListPrizes(activeOnly)
}

func (c *Catalog) Deactivate(id string) (model.PrizeDescription, error) {
	prize, err := c.Get(id)
	if err != nil {
		return model.PrizeDescription{}, err
	}
	prize.Active = false
	return prize, c.db.SavePrize(prize)
}

func (c *Catalog) Claim(id string) (model.PrizeDescription, error) {
	prize, err := c.Get(id)
	if err != nil {
		return model.PrizeDescription{}, err
	}
	if err := ValidateClaim(prize); err != nil {
		return model.PrizeDescription{}, err
	}
	prize.Claimed++
	return prize, c.db.SavePrize(prize)
}

func (c *Catalog) ClaimForRank(rank int, preferredTier string) (model.PrizeDescription, error) {
	prizes, err := c.List(true)
	if err != nil {
		return model.PrizeDescription{}, err
	}
	tier := NormalizeTier(preferredTier)
	if tier == "" {
		if selected, ok := SelectBest(prizes, rank); ok {
			return c.Claim(selected.ID)
		}
	} else {
		for _, prize := range Available(prizes, tier) {
			return c.Claim(prize.ID)
		}
	}
	return model.PrizeDescription{}, model.ErrPrizeExhausted
}
