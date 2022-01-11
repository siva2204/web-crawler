package crawler

import (
	"fmt"
	"sync"
	"time"

	"github.com/siva2204/web-crawler/config"
	neo4j_ "github.com/siva2204/web-crawler/neo4j"
	"github.com/siva2204/web-crawler/queue"
	redis_crawler "github.com/siva2204/web-crawler/redis"

	"github.com/siva2204/web-crawler/pagerank"
)

type HashMap struct {
	Hm map[string]bool
	sync.Mutex
}

// crawler bot type
type Crawler struct {
	Wg             sync.WaitGroup
	Threads        int
	Queue          *queue.Queue
	Hm             HashMap
	SeederListener chan int
	IsPaused       bool
	Ch             chan string
}

// run method starts crawling
func (c *Crawler) Run(graph *pagerank.PageRank, urlsRepository *neo4j_.Neo4jRepository) {
	// check if the url is already crawled

	if c.Queue.Len() == 0 {
		fmt.Println("queue is empty add some seed url to crawl")
		return
	}

	// ch := make(chan string, 10)

	for i := 0; i < c.Threads; i++ {
		c.Wg.Add(1)

		go func() {
			for {
				url, ok := <-c.Ch

				if !ok {
					c.Wg.Done()
					return
				}

				c.Hm.Lock()
				_, ok = c.Hm.Hm[url]
				c.Hm.Unlock()
				if !ok {
					fmt.Printf("crawling the %s url, now..", url)
					fmt.Println()

					// crawl with the url
					urls, err := uRLScrape(url, graph, urlsRepository)

					if err != nil {
						fmt.Printf("Error crawling url %+v", err)
						fmt.Println()

						// c.Wg.Done()
						// return
					}

					go func(url string) {
						data, err := dataScrape(url)

						if err != nil {
							fmt.Printf("Error getting data %+v", err)
							fmt.Println()

							// c.Wg.Done()
							// return
						}

						// for each token in data
						for _, token := range data {
							redis_crawler.Client.Append(token, url)
							urlsRepository.CreateToken(token)
							urlsRepository.ConnectTokenAndUrl(token, url)
						}

						// go rootNode.Insert(data, url)

					}(url)

					go func(url string) {
						description, tille, err := MetaScrape(url)

						if err != nil {
							fmt.Printf("Error getting description %+v", err)
							fmt.Println()
						}
						urlsRepository.AddDescriptionAndTitle(description, tille, url)
					}(url)

					c.Hm.Lock()
					c.Hm.Hm[url] = true
					c.Hm.Unlock()

					// enqueue the all the related url
					for _, urlc := range urls {
						c.Hm.Lock()

						_, ok := c.Hm.Hm[urlc]

						if !ok {
							c.Queue.Enqueue(urlc)
						}

						c.Hm.Unlock()
					}
				}
			}
		}()
	}

	// isPaused  := false

	go c.ListenForQueue()

	for {
		select {
		case status := <-c.SeederListener:
			fmt.Println("got a message : ", status)
			switch status {
			case 0:
				fmt.Println("Continue")
				// continue
				c.IsPaused = false
				go c.ListenForQueue()
				break
			case 1:
				fmt.Println("Pausing")
				// Pause
				c.IsPaused = true
				break
			}
		}
	}

	// traversing the queue
	// BFS

	close(c.Ch)
	c.Wg.Wait()
}

func (c *Crawler) ListenForQueue() {
	for {
		if c.IsPaused {
			fmt.Println("pausing !!!")
			break
		}
		if c.Queue.Len() != 0 {
			c.Ch <- c.Queue.Dequeue()
		}

		time.Sleep(time.Millisecond * time.Duration(config.Config.DequeueDelay))
	}
}
