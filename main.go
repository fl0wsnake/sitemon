package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	// "github.com/gocolly/colly/v2"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type ConfigEntry struct {
	Site   string
	Sel    string
	Filter string
	Emails []string
}

// `[{"site": "foo"}]`
func main() {
	var config_str = os.Getenv("SITEMON_CONFIG")
	println(config_str)
	var config []ConfigEntry
	err := json.Unmarshal([]byte(config_str), &config)
	if err != nil {
		println(err)
		return
	}

	for _, configEntry := range config {

		response, err := http.Get(configEntry.Site)
		if err != nil {
			log.Fatal(err)
		}
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(doc.Find(configEntry.Sel).Text())

		// c := colly.NewCollector()
		// c.OnHTML(configEntry.Sel,
		// 	func(e *colly.HTMLElement) {
		// 		fmt.Println(e.Text)
		// 	})
		// c.OnResponse(func(r *colly.Response) {
		// 	title := string(r.Body)
		// 	fmt.Println("title: " + title)
		// })
		// c.Visit(configEntry.Site)

	}
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

// // TODO: cloud function init
// func main() {
//    functions.HTTP("HelloHTTP", helloHTTP)
// }

// // helloHTTP is an HTTP Cloud Function with a request parameter.
// func helloHTTP(w http.ResponseWriter, r *http.Request) {
//   var d struct {
//     Name string `json:"name"`
//   }
//   if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
//     fmt.Fprint(w, "Hello, World!")
//     return
//   }
//   if d.Name == "" {
//     fmt.Fprint(w, "Hello, World!")
//     return
//   }
//   fmt.Fprintf(w, "Hello, %s!", html.EscapeString(d.Name))
// }
