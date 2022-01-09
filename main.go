package main

import (
	"flag"

	"github.com/siva2204/web-crawler/config"
	"github.com/siva2204/web-crawler/crawler"
	"github.com/siva2204/web-crawler/httpapi"
	"github.com/siva2204/web-crawler/queue"
	redis_crawler "github.com/siva2204/web-crawler/redis"
	"github.com/siva2204/web-crawler/trie"
)

var threads = flag.Int("threads", 100, "number of crawler threads")

func main() {
	// fmt.Println("work in progress")
	redis_crawler.CreateClient(config.Getenv("REDIS_HOST"), config.Getenv("REDIS_PORT"))
	// redis_crawler.Client.Insert("hello", []string{"a", "b", "c"})
	// redis_crawler.Client.Append("world", []string{"a", "b", "c"})
	// redis_crawler.CreateClient(config.Getenv("REDIS_HOST"), config.Getenv("REDIS_PORT"))

	crawler := crawler.Crawler{
		Threads: *threads,
		Queue:   &queue.Queue{},
		Hm: crawler.HashMap{
			Hm: make(map[string]bool),
		},
	}

	rootNode := trie.NewNode()
	crawler.Queue.Enqueue(config.Getenv("SEED_URL"))

	go crawler.Run(rootNode)

	httpapi.HttpServer()
}
