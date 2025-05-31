package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	// "testing"

	"cloud.google.com/go/storage"
	"github.com/PuerkitoBio/goquery"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Site_entry struct {
	Url       string
	Headers   map[string][]string
	Selectors []string
	Emails    []string
	Data      string
}

// var t = &testing.T{}

const bucket_name = "sitemon"
const object_name = "data"
const function_name = "sitemon"

func main() {
	ctx := context.Background()

	// TODO where do the credentials come from?
	client, err := storage.NewClient(ctx)
	assert(err, "Creating a GCS client")
	defer client.Close()

	object := client.Bucket(bucket_name).Object(object_name)
	sites_str := gcs_read(ctx, object)
	var sites []Site_entry
	err = json.Unmarshal([]byte(sites_str), &sites)
	assert(err)
	var mailing map[string]string
	for _, site := range sites {
		req, err := http.NewRequest("GET", site.Url, nil)
		req.Header = site.Headers
		client := http.Client{}
		res, err := client.Do(req)
		assert(err, "Fetching", site.Url)
		doc, err := goquery.NewDocumentFromReader(res.Body)
		assert(err, "Parsing response from", site.Url)
		var data_latest string
		for _, selector := range site.Selectors {
			data_latest += doc.Find(selector).Text()
		}

		if data_latest != site.Data {
			// TODO group sites by email address before sending
			for _, email := range site.Emails {
				mailing[email] += data_latest
			}
			gcs_write(ctx, object, data_latest)
		}

		for email, data := range mailing {
			siteurl, err := url.Parse(site.Url)
			assert(err)
			sendEmail(email, siteurl.Hostname(), data, site.Url)
		}
	}
}

func gcs_read(ctx context.Context, object *storage.ObjectHandle) string {
	reader, err := object.NewReader(ctx)
	assert(err)
	var contents []byte
	_, err = reader.Read(contents)
	assert(err, "Reading from gcs bucket")
	defer reader.Close()
	log.Printf("Blob %s downloaded: %s.\n", object.ObjectName(), contents)
	return string(contents)
}

func gcs_write(ctx context.Context, object *storage.ObjectHandle, contents string) {
	writer := object.NewWriter(ctx)
	_, err := io.WriteString(writer, contents)
	assert(err, "Writing to GCS bucket")
	defer writer.Close()
	log.Printf("Blob %v uploaded.\n", object.ObjectName())
}

// TODO link? Fix html content
func sendEmail(to, subject, content, link string) {
	fromEmail := mail.NewEmail("Tira", "tiramisu@example.com")
	toEmail := mail.NewEmail(to, to)
	htmlContent := fmt.Sprintf(`%s<a href="%s">`, content, link)
	message := mail.NewSingleEmail(fromEmail, subject, toEmail, content, htmlContent)
	// TODO sendgrid
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	assert(err)
	fmt.Println(response.StatusCode)
	fmt.Println(response.Headers)
}

func assert(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			log.Fatal(msg[0], ": ", err)
		} else {
			log.Fatal(err)
		}
	}
}
