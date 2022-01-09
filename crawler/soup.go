package crawler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/bbalet/stopwords"
	"github.com/jdkato/prose/v2"
)

var IsLetter = regexp.MustCompile(`^[a-z]+$`).MatchString

func uRLScrape(url string) ([]string, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return []string{}, fmt.Errorf("error status code %d", res.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return []string{}, err
	}

	// array of url
	var urls []string

	// Find the review items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the href
		href, _ := s.Attr("href")
		// push url to array
		urls = append(urls, href)
	})
	return urls, nil
}

func dataScrape(url string) (
	[]string, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// array of url
	var wordArray []string

	// Find the review items
	// doc.Find("p").Each(func(i int , s *goquery.Selection) {
	// For each item found, get the text
	text := doc.Text() // s.Text()

	//Return a string where HTML tags and French stop words has been removed
	cleanContent := stopwords.CleanString(text, "en", true)

	data, err := prose.NewDocument(cleanContent)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the doc's tokens:
	for _, tok := range data.Tokens() {
		// log.Println(tok.Text, tok.Tag, tok.Label)
		if IsLetter(tok.Text) {
			wordArray = append(wordArray, tok.Text)
		}
	}
	// })
	return wordArray, nil
}
