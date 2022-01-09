package main

import (
	"regexp"
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"github.com/jdkato/prose/v2"
	"github.com/bbalet/stopwords"
)

var IsLetter = regexp.MustCompile(`^[a-z]+$`).MatchString

func URLScrape(url string) ([]string, error) {
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
	var urls []string


	// Find the review items
	doc.Find("a").Each(func(i int , s *goquery.Selection) {
		// For each item found, get the href
		href, _ := s.Attr("href")
		// push url to array
		urls = append(urls, href)
	})
	return urls, nil
}

func PDataScrape(url string) (
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

func main(){
	urls, err := PDataScrape("http://localhost:5000/compellingly-embrace-from-generation-x-is")
	if err != nil {
		log.Fatal(err)
	}
	for _, url := range urls {
		log.Println(url)
	}
}
