package crawler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bbalet/stopwords"
	"github.com/jdkato/prose/v2"
	"github.com/siva2204/web-crawler/config"
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
		return []string{}, fmt.Errorf("error status code %d %s", res.StatusCode, url)
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
		// check if url has domain name
		if href != "" && href[0] == '/' {
			href = "https://" + strings.Split(url, "/")[2] + href
			urls = append(urls, href)
		} else if href != "" && href[0] != '/' && strings.Contains(href, config.Getenv("SEED_URL")) {
			// push url to array
			urls = append(urls, href)
		}
	})
	return urls, nil
}

func dataScrape(url string) (
	[]string, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return []string{}, fmt.Errorf("error status code %d %s", res.StatusCode, url)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return []string{}, err
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
		return []string{}, err
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

func MetaScrape(url string) (string, string, error) {

	var description string
	var title string

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", "", fmt.Errorf("error status code %d %s", res.StatusCode, url)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", "", err
	}

	text := doc.Text() // s.Text()

	cleanContent := stopwords.CleanString(text, "en", true)

	description = cleanContent[0:100]

	// array of url

	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		title = s.Text()
		fmt.Printf("Title field: %s\n", title)
	})

	return title, description, nil
}
