package main

import (
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

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

// func main(){
// 	urls, err := URLScrape("https://www.google.com/search?q=golang&tbm=isch")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, url := range urls {
// 		log.Println(url)
// 	}
// }