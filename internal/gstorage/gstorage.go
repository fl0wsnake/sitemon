package gstorage

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"sitemon/internal/monitoring"
	"sitemon/internal/util"

	"cloud.google.com/go/storage"
)

const bucket_name = "sitemon"
const object_name = "sites"

type GStorage struct {
	client *storage.Client
	object *storage.ObjectHandle
}

func Init(ctx context.Context) GStorage {
	client, err := storage.NewClient(ctx)
	util.Assert(err, "Creating a GS client")
	object := client.Bucket(bucket_name).Object(object_name)
	return GStorage{client, object}
}

func (gs *GStorage) Deinit() {
	gs.client.Close()
}

func (gs *GStorage) FetchSites(ctx context.Context) ([]monitoring.Site, error) {
	reader, err := gs.object.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	sites_byte, err := io.ReadAll(reader)
	var sites []monitoring.Site
	json.Unmarshal(sites_byte, &sites)
	util.Assert(err, "Reading from gs bucket")

	log.Printf("Blob %s downloaded: %s.\n", gs.object.ObjectName(), sites)
	return sites, nil
}

func (gs *GStorage) SaveSites(ctx context.Context, sites []monitoring.Site) {
	sites_byte, err := json.Marshal(sites)
	util.Assert(err, sites)

	writer := gs.object.NewWriter(ctx)
	_, err = io.Writer.Write(writer, sites_byte)
	util.Assert(err, "Writing to GS bucket")
	defer writer.Close()

	log.Printf("Blob %s uploaded to %v.\n", sites_byte, gs.object.ObjectName())
}
