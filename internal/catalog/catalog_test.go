package catalog

import (
	"path/filepath"
	"testing"

	"arcadeledger/internal/model"
	"arcadeledger/internal/store"
)

func TestCatalogClaimsInventory(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "arcade.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	c := New(db)
	if _, err := c.Add("pr1", "Champion Pin", "champion", "enamel pin", 1); err != nil {
		t.Fatal(err)
	}
	claimed, err := c.ClaimForRank(1, "")
	if err != nil {
		t.Fatal(err)
	}
	if claimed.Claimed != 1 || claimed.Tier != "champion" {
		t.Fatalf("unexpected claimed prize %#v", claimed)
	}
	if _, err := c.Claim("pr1"); err != model.ErrPrizeExhausted {
		t.Fatalf("expected exhaustion, got %v", err)
	}
}
