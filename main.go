package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sitemon/util"

	"cloud.google.com/go/storage"
	"github.com/PuerkitoBio/goquery"
)

type Site struct {
	Url       string
	Selectors []string
	Data      string
}

const bucket_name = "sitemon"
const object_name = "sites"

func main() {
	ctx := context.Background()

	// TODO where do the credentials come from?
	client, err := storage.NewClient(ctx)
	util.Assert(err, "Creating a GS client")
	defer client.Close()

	object := client.Bucket(bucket_name).Object(object_name)
	sites_str := gs_read(ctx, object)

	var sites []Site
	err = json.Unmarshal([]byte(sites_str), &sites)
	util.Assert(err)

	var notification_text string

	for _, site := range sites {
		req, err := http.NewRequest("GET", site.Url, nil)
		client := http.Client{}
		res, err := client.Do(req)
		util.Assert(err, "Fetching", site.Url)

		doc, err := goquery.NewDocumentFromReader(res.Body)
		util.Assert(err, "Parsing response from", site.Url)

		var data_latest string
		for _, selector := range site.Selectors {
			data_latest += doc.Find(selector).Text()
		}

		if data_latest != site.Data {
			notification_text += site.Url + ":\n" + data_latest + "\n\n"
			gs_write(ctx, object, data_latest)
		}
		fmt.Printf("NOTIFY_TRIGGER: %s", notification_text)
	}
}

func gs_read(ctx context.Context, object *storage.ObjectHandle) string {
	reader, err := object.NewReader(ctx)
	util.Assert(err)
	var contents []byte
	_, err = reader.Read(contents)
	util.Assert(err, "Reading from gs bucket")
	defer reader.Close()
	log.Printf("Blob %s downloaded: %s.\n", object.ObjectName(), contents)
	return string(contents)
}

func gs_write(ctx context.Context, object *storage.ObjectHandle, contents string) {
	writer := object.NewWriter(ctx)
	_, err := io.WriteString(writer, contents)
	util.Assert(err, "Writing to GS bucket")
	defer writer.Close()
	log.Printf("Blob %v uploaded.\n", object.ObjectName())
}
