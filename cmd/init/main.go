package sitemon

// For now only meant to be run locally once

import (
	"context"
	"encoding/json"

	"github.com/fl0wsnake/sitemon/internal/gstorage"
	"github.com/fl0wsnake/sitemon/internal/monitoring"
	"github.com/fl0wsnake/sitemon/internal/util"
	"os"
)

func main() error {
	ctx := context.Background

	gs := gstorage.Init(ctx)
	defer gs.Deinit()

	sites_byte, err := os.ReadFile("./sites.json")
	util.Assert(err)

	var sites []monitoring.Site
	err = json.Unmarshal(sites_byte, &sites)
	util.Assert(err)

	gs.SaveSites(ctx, sites)

	return nil
}
