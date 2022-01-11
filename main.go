package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/siva2204/web-crawler/config"
	"github.com/siva2204/web-crawler/crawler"
	"github.com/siva2204/web-crawler/db"
	"github.com/siva2204/web-crawler/httpapi"
	"github.com/siva2204/web-crawler/pagerank"
	"github.com/siva2204/web-crawler/queue"
	redis_crawler "github.com/siva2204/web-crawler/redis"
	"github.com/siva2204/web-crawler/trie"
	"github.com/siva2204/web-crawler/urls"
)

var threads = flag.Int("threads", 2, "number of crawler threads")

func main() {
	config.InitConfig()
	fmt.Printf("Initializing server with %d threads\n", *threads)
	redis_crawler.CreateClient(config.Config.RedisHost, config.Config.RedisPort)
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

	neo4jUri, found := os.LookupEnv("NEO4J_URI")
	if !found {
		panic("NEO4J_URI not set")
	}
	neo4jUsername, found := os.LookupEnv("NEO4J_USERNAME")
	if !found {
		panic("NEO4J_USERNAME not set")
	}
	neo4jPassword, found := os.LookupEnv("NEO4J_PASSWORD")
	if !found {
		panic("NEO4J_PASSWORD not set")
	}

	urlsRepository := urls.Neo4jRepository{
		Driver: driver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, "")),
	}

	urlsRepository.Init()

	urlsRepository.CreateUrl("https://www.google.com/")
	urlsRepository.CreateUrl("https://www.bing.com/")

	urlsRepository.AddPageRank(&urls.URL{
		URL:  "https://www.google.com/",
		RANK: 0.5,
	})

	rank, _ := urlsRepository.GetPageRank("https://www.google.com/")

	urlsRepository.ConnectTwoUrls("https://www.google.com/", "https://www.bing.com/")

	urls, _ := urlsRepository.GetConnectedUrls("https://www.google.com/")

	for _, url := range urls {
		fmt.Println(url.URL)
	}

	urlsRepository.CreateToken("token")

	urlsRepository.ConnectTokenAndUrl("token", "https://www.google.com/")

	fmt.Printf("Page rank: %f\n", rank)

	graph := pagerank.New()

	rootNode := trie.NewNode()
	crawlerBot.Queue.Enqueue(config.Config.SeedUrl)

	go crawlerBot.Run(graph)

	go crawler.SeederInstance.Run()

	httpapi.HttpServer(rootNode, graph)
}

func driver(target string, token neo4j.AuthToken) neo4j.Driver {
	result, err := neo4j.NewDriver(target, token)
	if err != nil {
		panic(err)
	}
	return result
}
