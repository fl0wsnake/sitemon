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
	"testing"

	"cloud.google.com/go/storage"
	"github.com/PuerkitoBio/goquery"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Site struct {
	Url       string
	Headers   map[string][]string
	Selectors []string
	Emails    []string
}

var t = &testing.T{}

const bucket_name = "sitemon"
const object_name = "data"
const function_name = "sitemon"

func main() {
	ctx := context.Background()
	var config_str = os.Getenv("CONFIG")
	var config []Site
	err := json.Unmarshal([]byte(config_str), &config)
	assert(err)

	// TODO where do the credentials come from?
	client, err := storage.NewClient(ctx)
	assert(err, "Creating a GCS client")
	defer client.Close()

	bucket := client.Bucket(bucket_name)
	for _, site := range config {
		req, err := http.NewRequest("GET", site.Url, nil)
		req.Header = site.Headers
		client := http.Client{}
		res, err := client.Do(req)
		assert(err, "Fetching", site.Url)
		doc, err := goquery.NewDocumentFromReader(res.Body)
		assert(err, "Parsing response from", site.Url)
		var site_text string
		for _, selector := range site.Selectors {
			site_text += doc.Find(selector).Text()
		}

		object := bucket.Object(object_name)
		bucket_text := gcs_read(ctx, object)
		if site_text != bucket_text {
			gcs_write(ctx, object, site_text)

			// Send mail
			for _, email := range site.Emails {
				siteurl, err := url.Parse(site.Url)
				assert(err)
				sendEmail(email, siteurl.Hostname(), site_text, site.Url)
			}
		}
	}
	return nil
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

func sendEmail(to, subject, content, link string) {
	fromEmail := mail.NewEmail("Tira", "tiramisu@example.com")
	toEmail := mail.NewEmail(to, to)
	htmlContent := fmt.Sprintf(`%s<a href="%s">`, content, link)
	message := mail.NewSingleEmail(fromEmail, subject, toEmail, content, htmlContent)
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
