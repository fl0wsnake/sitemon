package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"github.com/PuerkitoBio/goquery"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type ConfigEntry struct {
	Site     string
	Selector string
	Emails   []string
}

// `[{"site": "foo"}]`
func main() {
	var config_str = os.Getenv("SITEMON_CONFIG")
	var config []ConfigEntry
	err := json.Unmarshal([]byte(config_str), &config)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Config object: ", config)

	for _, configEntry := range config {
		response, err := http.Get(configEntry.Site)
		if err != nil {
			log.Fatal("Couldn't fetch", err)
		}
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal("Couldn't parse response", err)
		}
		text := doc.Find(configEntry.Selector).Text()
		if text == "" {
			log.Fatal("Selector returned an empty string: ", configEntry.Selector, "\n on site", configEntry.Site)
		}
		println("text:", text)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("Couldn't create a GCS client", err)
	}
	defer client.Close()

	object := client.Bucket(os.Getenv("bucket_name")).Object(os.Getenv("object_name"))
	gcs_read(ctx, object)
	// TODO
}

func gcs_read(ctx context.Context, object *storage.ObjectHandle) string {
	reader, err := object.NewReader(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var contents []byte
	_, err = reader.Read(contents)
	if err != nil {
		log.Fatal("Error reading from GCS bucket", err)
	}
	defer reader.Close()
	log.Printf("Blob %s downloaded: %s.\n", object.ObjectName(), contents)
	return string(contents)
}

func gcs_write(ctx context.Context, object *storage.ObjectHandle, contents string) {
	writer := object.NewWriter(ctx)
	_, err := io.WriteString(writer, contents)
	if err != nil {
		log.Fatal("Error writing to GCS bucket", err)
	}
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
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Headers)
	}
}
