package monitoring

import (
	"fmt"
	"net/http"
	"regexp"
	"sitemon/internal/util"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Site struct {
	Url         string
	Matchers    map[string]string
	MatchedText string
}

func (site *Site) Match() *string {
	req, err := http.NewRequest("GET", site.Url, nil)
	client := http.Client{}
	res, err := client.Do(req)
	util.Assert(err, "Fetching", site.Url)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	util.Assert(err, "Parsing response from", site.Url)

	var text_matched_sb strings.Builder
	for selector, regex := range site.Matchers {
		text_selected := doc.Find(selector).Text()
		did_match, err := regexp.MatchString(regex, text_selected)
		util.Assert(err, "Compiling regex", regex)

		if did_match {
			text_matched_sb.WriteString(text_selected)
		}
	}

	if text_matched_sb.String() != site.MatchedText {
		site.MatchedText = text_matched_sb.String() // BUG: does not seem to work
		res := fmt.Sprintf("# %s:\n%s\n\n", site.Url, text_matched_sb.String())
		return &res
	} else {
		return nil
	}
}
