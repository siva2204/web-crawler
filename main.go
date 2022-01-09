package main

import (
	"fmt"

	redis_crawler "github.com/siva2204/web-crawler/redis"
)

func main() {
	fmt.Println("work in progress")
	redis_crawler.CreateClient(Getenv("REDIS_HOST"), Getenv("REDIS_PORT"))
	redis_crawler.Client.Insert("hello", []string{"a", "b", "c"})
	redis_crawler.Client.Append("hello", []string{"a", "b", "c"})
}
