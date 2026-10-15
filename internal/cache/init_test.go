package cache_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/mia0x75/halo/internal/app"
	"github.com/mia0x75/halo/internal/cache"
	"github.com/mia0x75/halo/internal/platform/db"
)

func TestInitSmoke(t *testing.T) {
	data, err := os.ReadFile("../../bin/config.json")
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

	cache.Init(application.Repos)
	defer cache.LoopInit()

	t.Logf("users=%d rules=%d clusters=%d options=%d edges=%d avatars=%d",
		cache.UsersMap.Count(), cache.RulesMap.Count(), cache.ClustersMap.Count(),
		cache.OptionsMap.Count(), cache.EdgesMap.Count(), cache.AvatarsMap.Count())
	if cache.UsersMap.Count() == 0 || cache.RulesMap.Count() == 0 {
		t.Fatal("caches did not load")
	}
}
