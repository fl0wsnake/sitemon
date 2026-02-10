package main

import (
	"context"
	"fmt"
	"sitemon/internal/gstorage"
	"sitemon/internal/util"
	"strings"
)

func main() {
	ctx := context.Background()
	gs := gstorage.Init(ctx)
	defer gs.Deinit()

	sites, err := gs.FetchSites(ctx)
	util.Assert(err, "GS object does not exist")

	var matches strings.Builder
	for _, site := range sites {
		match := site.Match() // TODO parallelize
		if match != nil {
			matches.WriteString(*match)
		}
	}
	if matches.Len() > 0 {
		gs.SaveSites(ctx, sites) // TODO parallelize
		fmt.Printf("ALERT_TRIGGER: %s", matches.String())
	}
}
