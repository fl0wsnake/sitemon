package sitemon

import (
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/fl0wsnake/sitemon/internal/gstorage"
	"github.com/fl0wsnake/sitemon/internal/util"
	"strings"
)

func init() {
	functions.CloudEvent("HourlyFunc", hourlyFunc)
}

func hourlyFunc(ctx context.Context, e event.Event) error {
	gs := gstorage.Init(ctx)
	defer gs.Deinit()

	sites, err := gs.FetchSites(ctx)
	util.Assert(err, "GS object does not exist")

	var matches strings.Builder
	for site_i := range sites {
		match := sites[site_i].Match() // TODO parallelize
		if match != nil {
			matches.WriteString(*match)
		}
	}
	if matches.Len() > 0 {
		gs.SaveSites(ctx, sites) // TODO parallelize
		fmt.Printf("ALERT_TRIGGER: %s", matches.String())
	}

	return nil
}
