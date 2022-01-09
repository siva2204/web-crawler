package main

import (
	"github.com/siva2204/web-crawler/crawler"
	"github.com/siva2204/web-crawler/queue"
	redis_crawler "github.com/siva2204/web-crawler/redis"
	"github.com/siva2204/web-crawler/config"
)

func main() {
	// fmt.Println("work in progress")
	redis_crawler.CreateClient(config.Getenv("REDIS_HOST"), config.Getenv("REDIS_PORT"))
	// redis_crawler.Client.Insert("hello", []string{"a", "b", "c"})
	// redis_crawler.Client.Append("world", []string{"a", "b", "c"})

	crawler := crawler.Crawler{
		Threads: 100,
		Queue:   &queue.Queue{},
		Hm: crawler.HashMap{
			Hm: make(map[string]bool),
		},
	}

	crawler.Queue.Enqueue("http://localhost:5000")

	crawler.Run()

}
