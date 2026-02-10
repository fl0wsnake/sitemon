package main

// For now only meant to be run locally once

import (
	"context"
	"encoding/json"
	"os"
	"sitemon/internal/gstorage"
	"sitemon/internal/monitoring"
	"sitemon/internal/util"
)

func main() {
	ctx := context.Background()
	gs := gstorage.Init(ctx)
	defer gs.Deinit()

	sites_byte, err := os.ReadFile("./sites.json")
	util.Assert(err)

	var sites []monitoring.Site
	err = json.Unmarshal(sites_byte, &sites)
	util.Assert(err)

	gs.SaveSites(ctx, sites)
}
