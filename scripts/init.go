package init

// For now only meant to be run locally once

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"sitemon/util"

	// "testing"

	"cloud.google.com/go/storage"
)

func main() {
	ctx := context.Background()

	// TODO where do the credentials come from?
	client, err := storage.NewClient(ctx)
	util.Assert(err, "Creating a GCS client")
	defer client.Close()

	object := client.Bucket(bucket_name).Object(object_name)
	reader, err := object.NewReader(ctx)
	if err == nil {
		var contents []byte
		_, err = reader.Read(contents)
		util.Assert(err, "Reading from gs bucket")
		defer reader.Close()

		log.Printf("Blob %s downloaded: %s.\n", object.ObjectName(), contents)
		sites_str := string(contents)

		var sites []Site
		err = json.Unmarshal([]byte(sites_str), &sites)
		util.Assert(err)

		log.Fatal(err, sites, "Sites object already exists")
	}

	sites_str, err := os.ReadFile("./sites.json")
	util.Assert(err)

	writer := object.NewWriter(ctx)
	_, err = io.WriteString(writer, string(sites_str))
	util.Assert(err, "Writing to GS bucket")
	defer writer.Close()

	log.Printf("Blob %v uploaded.\n", object.ObjectName())
}
