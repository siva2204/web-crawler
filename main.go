package main

import (
	"flag"

	"github.com/siva2204/web-crawler/config"
	"github.com/siva2204/web-crawler/crawler"
	"github.com/siva2204/web-crawler/db"
	"github.com/siva2204/web-crawler/httpapi"
	"github.com/siva2204/web-crawler/queue"
	redis_crawler "github.com/siva2204/web-crawler/redis"
	"github.com/siva2204/web-crawler/trie"
)

var threads = flag.Int("threads", 2, "number of crawler threads")

func main() {
	redis_crawler.CreateClient(config.Getenv("REDIS_HOST"), config.Getenv("REDIS_PORT"))
	db.InitDB()

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

	rootNode := trie.NewNode()
	crawlerBot.Queue.Enqueue(config.Getenv("SEED_URL"))

	go crawlerBot.Run(rootNode)

	go crawler.SeederInstance.Run()

	crawlerBot.Queue.Enqueue(config.Getenv("SEED_URL"))

	httpapi.HttpServer(rootNode)
}
