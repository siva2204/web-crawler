package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/siva2204/web-crawler/config"
	"github.com/siva2204/web-crawler/crawler"
	"github.com/siva2204/web-crawler/queue"
	redis_crawler "github.com/siva2204/web-crawler/redis"
)

var threads = flag.Int("threads", 2, "number of crawler threads")

func main() {
	// fmt.Println("work in progress")
	redis_crawler.CreateClient(config.Getenv("REDIS_HOST"), config.Getenv("REDIS_PORT"))
	// redis_crawler.Client.Insert("hello", []string{"a", "b", "c"})
	// redis_crawler.Client.Append("world", []string{"a", "b", "c"})
	// redis_crawler.CreateClient(config.Getenv("REDIS_HOST"), config.Getenv("REDIS_PORT"))

	crawlerBot := crawler.Crawler{
		Threads: *threads,
		Queue:   &queue.Queue{},
		Hm: crawler.HashMap{
			Hm: make(map[string]bool),
		},
		IsPaused: false,
		Ch:       make(chan string, 50),
	}
	crawler.InitSeeder(&crawlerBot)

	go crawlerBot.Run()

	go crawler.SeederInstance.Run()

	crawlerBot.Queue.Enqueue("http://localhost:5000")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		}
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			urls, err := redis_crawler.Client.GetUnEncoded("hello")
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"status": "error"})
			} else {
				json.NewEncoder(w).Encode(urls)
			}
		}
	})

	http.ListenAndServe(":7000", nil)
}
