package caches_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/mia0x75/halo/caches"
	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/internal/platform/db"
)

func TestInitSmoke(t *testing.T) {
	data, err := os.ReadFile("../bin/config.json")
	if err != nil {
		t.Skip("no bin/config.json")
	}
	var cfg struct {
		Database struct {
			Addr string `json:"addr"`
		} `json:"database"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Database.Addr == "" {
		t.Skip("no db addr")
	}

	application, err := app.New(context.Background(), db.Config{Addr: cfg.Database.Addr})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = application.Close() }()

	caches.Init(application.Repos)
	defer caches.LoopInit()

	t.Logf("users=%d rules=%d clusters=%d options=%d edges=%d avatars=%d",
		caches.UsersMap.Count(), caches.RulesMap.Count(), caches.ClustersMap.Count(),
		caches.OptionsMap.Count(), caches.EdgesMap.Count(), caches.AvatarsMap.Count())
	if caches.UsersMap.Count() == 0 || caches.RulesMap.Count() == 0 {
		t.Fatal("caches did not load")
	}
}
