package main

import (
	"encoding/json"
	"os"
	"fmt"
	"github.com/gocolly/colly/v2"
)

type ConfigEntry struct {
	Site string
	// selector string
	// emails   []string
}

// `[{"site": "foo"}]`
func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("#BigProductCard > div.jsx-b91b1cc89b9c436e.BigProductCard__content > div:nth-child(2) > div.jsx-7d9b01de45811631.BigProductCardTopInfo > div.jsx-7d9b01de45811631.BigProductCardTopInfo__priceInfo > h6 > .Price__value_title", 
	func(e *colly.HTMLElement) {
		// txt := e.Attr("innerText")
		fmt.Println(e.Text)
	})

	c.Visit("https://metro.zakaz.ua/uk/products/ukrayina--04820254610291/")

	var config_str = os.Getenv("SITEMON_CONFIG")
	println(config_str)
	var config []ConfigEntry
	err := json.Unmarshal([]byte(config_str), &config)
	if err != nil {
		println(err)
		return
	}

	// curl_stdout, err := exec.Command("curl", "https://metro.zakaz.ua/uk/products/ukrayina--04820254610291/",
	// 	"-H", "accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
	// 	"-H", "accept-language: en-US,en;q=0.8",
	// 	"-H", "cache-control: max-age=0",
	// 	"-H", "cookie: __zlcmid=1MsmszNtNM1UzBW; storeId=48215637; deliveryType=plan",
	// 	"-H", "if-modified-since: Mon, 22 Jul 2024 00:53:00 GMT",
	// 	"-H", "priority: u=0, i",
	// 	"-H", `sec-ch-ua: "Not/A)Brand";v="8", "Chromium";v="126", "Brave";v="126"`,
	// 	"-H", "sec-ch-ua-mobile: ?0",
	// 	"-H", `sec-ch-ua-platform: "Linux"`,
	// 	"-H", "sec-fetch-dest: document",
	// 	"-H", "sec-fetch-mode: navigate",
	// 	"-H", "sec-fetch-site: same-origin",
	// 	"-H", "sec-fetch-user: ?1",
	// 	"-H", "sec-gpc: 1",
	// 	"-H", "upgrade-insecure-requests: 1",
	// 	"-H", "user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
	// ).Output()

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // fmt.Println(string(curl_stdout))

	// pup_cmd := exec.Command("pup", `#BigProductCard > div.jsx-b91b1cc89b9c436e.BigProductCard__content > div:nth-child(2) > div.jsx-7d9b01de45811631.BigProductCardTopInfo > div.jsx-7d9b01de45811631.BigProductCardTopInfo__priceInfo > h6 > .Price__value_title text{}`)
	// pup_pipe, err := pup_cmd.StdinPipe()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// n, err := pup_pipe.Write(curl_stdout)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(n)

	// pup_stdout, err := pup_cmd.Output()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(pup_stdout)
	// pup_pipe.Close()
}

// // TODO: init
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
